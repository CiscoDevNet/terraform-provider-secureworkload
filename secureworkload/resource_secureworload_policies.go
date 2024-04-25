package secureworkload

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
)

func resourceSecureWorkloadPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a new policy on Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_policies\" \"policy1\" {\n" +
			"	 workspace_id = secureworkload_workspace.workspace.id\n" +
			"	 consumer_filter_id = secureworkload_filter.any.id\n" +
			"	 provider_filter_id = secureworkload_cluster.web.id\n" +
			"    policy_action = \"ALLOW\"\n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.\n",
		Create: resourceSecureWorkloadPolicyCreate,
		Update: nil,
		Read:   resourceSecureWorkloadPolicyRead,
		Delete: resourceSecureWorkloadPolicyDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the needed workspace.",
			},
			"consumer_filter_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a defined filter.",
			},
			"provider_filter_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of a defined filter.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates the version of the workspace the cluster will be added to.",
			},
			"rank": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Values can be DEFAULT, ABSOLUTE or CATCHALL for ranking",
			},
			"policy_action": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Values can be ALLOW or DENY: means whether we should allow or drop traffic from consumer to provider on the given service port/protocol",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Used to sort policy.",
			},
			// "protocol": {
			// 	Type:        schema.TypeInt,
			// 	ForceNew:    true,
			// 	Optional:    true,
			// 	Description: "Protocol integer value (NULL means all protocols).",
			// },
			// "start_port": {
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// 	ForceNew:    true,
			// 	Description: "Start port of the range.",
			// },
			// "end_port": {
			// 	Type:        schema.TypeInt,
			// 	Optional:    true,
			// 	ForceNew:    true,
			// 	Description: "End port of the range.",
			// },
			// "policy_approved": {
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	ForceNew:    true,
			// 	Default:     false,
			// 	Description: "(Optional) Indicates whether the policy is approved. Default is false.",
			// },
		},
	}
}

var requiredCreatePolicyParams = []string{"consumer_filter_id", "provider_filter_id", "policy_action"}

func resourceSecureWorkloadPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreatePolicyParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	// proto := d.Get("protocol").(int)
	// ranges := [2]int{d.Get("start_port").(int), d.Get("end_port").(int)}
	// approved := d.Get("policy_approved").(bool)
	// L4Port := []Layer4Network{
	// 	{
	// 		Protocol:  proto,
	// 		PortRange: ranges,
	// 		Approved:  approved,
	// 	},
	// }
	createPolicyParams := CreatePolicyRequest{
		ConsumerId: d.Get("consumer_filter_id").(string),
		ProviderId: d.Get("provider_filter_id").(string),
		Version:    d.Get("version").(string),
		Rank:       d.Get("rank").(string),
		Action:     d.Get("policy_action").(string),
		Priority:   d.Get("priority").(int),
		// Layer4NetworkPolicies: L4Port,
	}
	policy, err := client.CreatePolicy(createPolicyParams, d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	d.SetId(policy.Id)
	return nil
}

func resourceSecureWorkloadPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	policy, err := client.DescribePolicy(d.Id())
	if err != nil {
		return err
	}
	d.Set("consumer_filter_id", policy.ConsumerId)
	d.Set("provider_filter_id", policy.ProviderId)
	d.Set("version", policy.Version)
	d.Set("rank", policy.Rank)
	d.Set("policy_action", policy.Action)
	d.Set("priority", policy.Priority)
	return nil
}

func resourceSecureWorkloadPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeletePolicy(d.Get("workspace_id").(string), d.Id())
}
