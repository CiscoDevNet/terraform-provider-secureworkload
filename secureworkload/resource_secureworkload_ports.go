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
			"**Note:** If creating multiple resources for ports during a single `terraform apply`, you may have to use `depends_on` to chain the resources so that terraform creates it in the same order that you intended.\n",
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

func resourceSecureWorkloadPortRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	policyId := d.Get("policy_id").(string)
	policy, err := client.DescribePolicyL4Params(policyId)
	if err != nil {
		return err
	}

	var found bool
	var port Port
	for i := range policy.L4Params {
		if policy.L4Params[i].Id == d.Id() {
			port = policy.L4Params[i]
			found = true
			break
		}
	}

	if !found {
		// The l4_param no longer exists (e.g. deleted out of band).
		// Signal to Terraform that the resource is gone so it can be
		// recreated on the next apply.
		d.SetId("")
		return nil
	}

	// The API returns the range as "port": [start, end] rather than
	// start_port/end_port, so backfill before writing to state.
	port.normalize()

	d.Set("start_port", port.StartPort)
	d.Set("end_port", port.EndPort)
	d.Set("description", port.Description)
	d.Set("proto", port.Proto)
	if port.Version != "" {
		d.Set("version", port.Version)
	}
	return nil
}

func resourceSecureWorkloadPortDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeletePort(d.Get("policy_id").(string), d.Id())
}
