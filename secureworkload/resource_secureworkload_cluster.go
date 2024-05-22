package secureworkload

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
)

func resourceSecureWorkloadCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a new Cluster on Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_cluster\" \"cluster\" {\n" +
			"	 workspace_id = data.secureworkload_workspace.workspace.id\n" +
			"    name = \"New-Cluster\"\n" +
			"    description = \"New Cluster via TF\"\n" +
			"    query = <<EOF\n" +
			"                {" +
			"        		 \"type\":\"eq\",\n" +
			"        		 \"field\": \"ip\",\n" +
			"        		 \"value\": \"10.0.0.1\"\n" +
			"        		 }\n" +
			"        	EOF\n" +
			"    approved = false \n" +
			"}\n" +
			"resource \"secureworkload_cluster\" \"cluster2\" {\n" +
			"    depends_on = [secureworkload_cluster.cluster] \n" +
			"	 workspace_id = data.secureworkload_workspace.workspace.id\n" +
			"    name = \"New-Cluster2\"\n" +
			"    description = \"Second Cluster via TF\"\n" +
			"    query = <<EOF\n" +
			"                {" +
			"        		 \"type\":\"eq\",\n" +
			"        		 \"field\": \"ip\",\n" +
			"        		 \"value\": \"10.0.0.2\"\n" +
			"        		 }\n" +
			"        	EOF\n" +
			"    approved = false \n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple clusters during a single `terraform apply`, remember to use `depends_on` to chain the filters so that terraform creates them in a specific order to avoid *429:too_many_request* error.\n",
		Create: resourceSecureWorkloadClusterCreate,
		Update: nil,
		Read:   resourceSecureWorkloadClusterRead,
		Delete: resourceSecureWorkloadClusterDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "User-specified name for the inventory cluster.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates the version of the workspace the cluster will be added to.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "(optional) The description of the cluster.",
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "JSON object representation of an inventory cluster query.",
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the workspace associated with the cluster.",
			},
			"approved": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "(optional) An approved cluster will not be updated during an automatic policy discovery run. Default false.",
			},
		},
	}
}

var requiredCreateClusterParams = []string{"name", "workspace_id", "query"}

func resourceSecureWorkloadClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreateClusterParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	createClusterParams := CreateClusterRequest{
		Name:        d.Get("name").(string),
		Version:     d.Get("version").(string),
		Description: d.Get("description").(string),
		Query:       []byte(d.Get("query").(string)),
		Approved:    d.Get("approved").(bool),
	}
	cluster, err := client.CreateCluster(createClusterParams, d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	d.SetId(cluster.Id)
	return nil
}

func resourceSecureWorkloadClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	cluster, err := client.DescribeCluster(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", cluster.Name)
	d.Set("version", cluster.Version)
	d.Set("description", cluster.Description)
	d.Set("approved", cluster.Approved)
	return nil
}

func resourceSecureWorkloadClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteCluster(d.Get("workspace_id").(string), d.Id())
}
