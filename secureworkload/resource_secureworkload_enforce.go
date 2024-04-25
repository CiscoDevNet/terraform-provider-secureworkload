package secureworkload

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
)

func resourceSecureWorkloadEnforce() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for enforcing policy on a single workspace.\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_enforce\" \"enforced\" {\n" +
			"	 workspace_id = secureworkload_workspace.workspace.id\n" +
			"    version = \"p10\" \n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.\n",
		Create: resourceSecureWorkloadEnforceCreate,
		Update: nil,
		Read:   resourceSecureWorkloadEnforceRead,
		Delete: resourceSecureWorkloadEnforceDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
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
		},
	}
}

var requiredCreateEnforceParams = []string{"workspace_id"}

func resourceSecureWorkloadEnforceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreateEnforceParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	createEnforceParams := CreateEnforceRequest{
		Version: d.Get("version").(string),
	}
	port, err := client.CreateEnforce(createEnforceParams, d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	d.SetId(port.Epoch)
	return nil
}

func resourceSecureWorkloadEnforceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	describeApplicatioParams := DescribeApplicationRequest{
		ApplicationId: d.Get("workspace_id").(string),
	}
	application, err := client.DescribeApplication(describeApplicatioParams)
	if err != nil {
		return err
	}
	d.Set("name", application.Name)
	d.Set("description", application.Description)
	d.Set("primary", application.Primary)
	d.Set("alternate_query_mode", application.AlternateQueryMode)
	d.Set("latest_adm_version", application.LatestADMVersion)
	d.Set("enforcement_enabled", application.EnforcementEnabled)
	d.Set("enforced_version", application.EnforcedVersion)
	return nil
}

func resourceSecureWorkloadEnforceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteEnforce(d.Get("workspace_id").(string))
}
