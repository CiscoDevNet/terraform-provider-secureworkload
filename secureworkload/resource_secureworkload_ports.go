package secureworkload

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
)

func resourceSecureWorkloadPort() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a new service port on Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_port\" \"port1\" {\n" +
			"	 policy_id = secureworkload_policies.policy1.id\n" +
			"    start_port = 80 \n" +
			"    end_port = 80 \n" +
			"    proto = 6 \n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.\n",
		Create: resourceSecureWorkloadPortCreate,
		Update: nil,
		Read:   resourceSecureWorkloadPortRead,
		Delete: resourceSecureWorkloadPortDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the needed policy.",
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
				Description: "(optional) Short string about this proto and port",
			},
			"start_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Start port of the range.",
			},
			"end_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "End port of the range.",
			},
			"proto": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Protocol Integer value (NULL means all protocols)",
			},
		},
	}
}

var requiredCreatePortParams = []string{"policy_id", "start_port", "end_port"}

func resourceSecureWorkloadPortCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreatePortParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	createPortParams := CreatePortRequest{
		StartPort:   d.Get("start_port").(int),
		EndPort:     d.Get("end_port").(int),
		Version:     d.Get("version").(string),
		Description: d.Get("description").(string),
		Proto:       d.Get("proto").(int),
	}
	port, err := client.CreatePort(createPortParams, d.Get("policy_id").(string))
	if err != nil {
		return err
	}
	d.SetId(port.Id)
	return nil
}

// func resourceSecureWorkloadPortRead(d *schema.ResourceData, meta interface{}) error {
// 	client := meta.(Client)
// 	port, err := client.DescribePort(d.Get("policy_id").(string), d.Id())
// 	if err != nil {
// 		return err
// 	}
// 	d.Set("start_port", port.StartPort)
// 	d.Set("end_port", port.EndPort)
// 	d.Set("version", port.Version)
// 	d.Set("description", port.Description)
// 	d.Set("proto", port.Proto)
// 	return nil
// }

func resourceSecureWorkloadPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	policy, err := client.DescribePolicy(d.Get("policy_id").(string))
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

func resourceSecureWorkloadPortDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeletePort(d.Get("policy_id").(string), d.Id())
}
