package secureworkload

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
	// secureworkload "github.com/secureworkload-exchange/terraform-go-sdk"
)

var timer int

func resourceSecureWorkloadScope() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a scope in Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_scope\" \"scope\" {\n" +
			"    short_name = \"Terraform-created-scope\"\n" +
			"    sub_type = \"DNS_SERVERS\"\n" +
			"    short_query = file(\"${path.module}/query_file.json\") \n" +
			"	 parent_app_scope_id = data.secureworkload_scope.scope.id\n" +
			"}\n" +
			"\n" +
			"resource \"secureworkload_scope\" \"scope2\" {\n" +
			"    short_name = \"Terraform-created-scope2\"\n" +
			"    query = <<EOF\n" +
			"                { \n" +
			"        		 \"type\":\"subnet\",\n" +
			"        		 \"field\": \"ip\",\n" +
			"        		 \"value\": \"10.0.1.0/24\"\n" +
			"        		 }\n" +
			"        	EOF\n" +
			"    sub_type = \"GENERIC\"\n" +
			"	 parent_app_scope_id = data.secureworkload_scope.scope.id\n" +
			"}\n" +
			"```\n" +
			"Example of **query_file.json**\n" +
			"```json\n" +
			"{	\n" +
			"	\"type\": \"or\",\n" +
			"	\"filters\": [ \n" +
			"	  {\n" +
			"	\"type\": \"and\",\n" +
			"		\"filters\": [ \n" +
			"		  { \n" +
			"			\"type\": \"contains\",\n" +
			"			\"field\": \"user_orchestrator_system/name\",\n" +
			"			\"value\": \"Random\"\n" +
			"		  },\n" +
			"		  {\n" +
			"			\"type\": \"eq\",\n" +
			"			\"field\": \"ip\",\n" +
			"			\"value\": \"10.0.1.1\"\n" +
			"		  }\n" +
			"		]\n" +
			"	  },\n" +
			"	  {\n" +
			"		\"type\": \"gt\",\n" +
			"		\"field\": \"host_tags_cvss\",\n" +
			"		\"value\": 2\n" +
			"	  }\n" +
			"	]\n" +
			"  }\n" +
			"```\n" +
			"**Note:** If creating multiple resources for scope during a single `terraform apply`, you may have to use `depends_on` to chain the resources so that terraform creates it in the same order that you intended.\n",
		Create: resourceSecureWorkloadScopeCreate,
		Update: resourceSecureWorkloadScopeUpdate,
		Read:   resourceSecureWorkloadScopeRead,
		Delete: resourceSecureWorkloadScopeDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"short_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User-specified name for the scope.",
			},
			"sub_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "User-specified sub type for the scope.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "User-specified description of the scope.",
			},
			"parent_app_scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "What resource field to use when evaluating the scope query.",
			},
			"policy_priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    false,
				Default:     nil,
				Computed:    true,
				Description: "Used to sort application priorities; default is last.",
			},
			"short_query": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "JSON object representation of an inventory filter query. The query shown in the above example is 'orchestrator_system/name containes Random and Address = 10.0.1.1 or CVE Score v3 >2'.Operator can any of the following: [and, or, eq, subnet, contains, regex, gt, gte, lt, lte, in, range, ranges, not, all, none] ",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fully qualified name of the scope. This is a fully qualified name; that is, it includes the names of parent scopes (if applicable) all the way to the root scope.",
			},
			"root_app_scope_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Root scope for the secureworkload installation",
			},
			"vrf_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the VRF to which scope belongs.",
			},
			"priority": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"short_priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Used to sort application priorities; default is last.",
			},
			"dirty": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates a child or parent query has been updated and that the changes need to be committed..",
			},
			"child_app_scope_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Indicates a child or parent query has been updated and that the changes need to be committed..",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix Epoch timestamp when scope was created.",
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Unix Epoch timestamp when scope was last updated.",
			},
		},
	}
}

var requiredCreateScopeParams = []string{"short_name", "parent_app_scope_id"}

func resourceSecureWorkloadScopeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreateScopeParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	createScopeParams := CreateScopeRequest{
		ShortName:        d.Get("short_name").(string),
		Description:      d.Get("description").(string),
		ParentAppScopeId: d.Get("parent_app_scope_id").(string),
		SubType:          d.Get("sub_type").(string),
		ShortQuery:       []byte(d.Get("short_query").(string)),
		PolicyPriority:   d.Get("policy_priority").(int),
	}
	scope, err := client.CreateScope(createScopeParams)
	if err != nil {
		return err
	}
	d.Set("policy_priority", scope.PolicyPriority)
	d.Set("description", scope.Description)
	d.SetId(scope.Id)
	return nil
}

func resourceSecureWorkloadScopeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	scope, err := client.DescribeScope(d.Id())
	if err != nil {
		return err
	}
	d.Set("short_name", scope.ShortName)
	d.Set("description", scope.Description)
	d.Set("parent_app_scope_id", scope.ParentAppScopeId)
	d.Set("policy_priority", scope.PolicyPriority)
	d.Set("name", scope.Name)
	d.Set("root_app_scope_id", scope.RootAppScopeId)
	d.Set("vrf_id", scope.VRFId)
	d.Set("priority", scope.Priority)
	d.Set("short_priority", scope.ShortPriority)
	d.Set("dirty", scope.Dirty)
	d.Set("child_app_scope_ids", scope.ChildAppScopeIds)
	d.Set("created_at", scope.CreatedAt)
	d.Set("updated_at", scope.UpdatedAt)
	return nil
}

func resourceSecureWorkloadScopeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	client.DeleteScope(d.Id())
	createScopeParams := CreateScopeRequest{
		ShortName:        d.Get("short_name").(string),
		Description:      d.Get("description").(string),
		ParentAppScopeId: d.Get("parent_app_scope_id").(string),
		ShortQuery:       []byte(d.Get("short_query").(string)),
		SubType:          d.Get("sub_type").(string),
		PolicyPriority:   d.Get("policy_priority").(int),
	}
	scope, err := client.CreateScope(createScopeParams)
	if err != nil {
		return err
	}
	d.Set("description", scope.Description)
	d.Set("policy_priority", scope.PolicyPriority)
	d.SetId(scope.Id)
	return nil
}

func resourceSecureWorkloadScopeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	err := client.DeleteScope(d.Id())
	for err != nil {
		if strings.Contains(err.Error(), "error:cannot delete scope because it is in use") {
			if timer >= 20 {
				return err
			}
			time.Sleep(60 * time.Second)
			timer = timer + 1
			err = client.DeleteScope(d.Id())
		} else {
			return err
		}
	}
	return err
}
