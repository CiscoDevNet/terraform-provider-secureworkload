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
			"To always enforce the newest available policy version on each `terraform apply` (without knowing the version ahead of time), omit `version` and set `track_latest_version = true`:\n" +
			"```hcl\n" +
			"resource \"secureworkload_enforce\" \"enforced\" {\n" +
			"	 workspace_id          = secureworkload_workspace.workspace.id\n" +
			"    track_latest_version = true\n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.\n",
		Create: resourceSecureWorkloadEnforceCreate,
		Update: resourceSecureWorkloadEnforceUpdate,
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "The policy version to enforce, in the form \"p10\". If omitted and " +
					"`track_latest_version` is `true`, the newest available policy version is enforced " +
					"and this attribute tracks it automatically on every `terraform apply`. Changing this " +
					"value now updates the enforcement in place instead of destroying and recreating the resource.",
			},
			"track_latest_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "When `true`, always enforce the newest available policy version for the " +
					"workspace on every `terraform apply`, without needing to know or input the version.",
			},
			"enforcement_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if enforcement is currently enabled on the workspace.",
			},
			"enforced_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy version currently enforced on the workspace, in the form \"p10\".",
			},
			"latest_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The newest available policy version for the workspace, in the form \"p10\".",
			},
		},
	}
}

var requiredCreateEnforceParams = []string{"workspace_id"}

// resolveEnforceVersion determines which policy version should be enforced based on
// the resource configuration: an explicit `version` wins, otherwise
// `track_latest_version` resolves to the newest available version, otherwise an
// empty version is sent (preserving the historical default API behaviour).
func resolveEnforceVersion(d *schema.ResourceData, client Client, workspaceId string) (string, error) {
	if v, ok := d.GetOk("version"); ok {
		if s, ok := v.(string); ok && s != "" {
			return s, nil
		}
	}
	if d.Get("track_latest_version").(bool) {
		return client.LatestPolicyVersion(workspaceId)
	}
	return "", nil
}

func resourceSecureWorkloadEnforceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	for _, param := range requiredCreateEnforceParams {
		if d.Get(param) == "" {
			return fmt.Errorf("%s is required but was not provided", param)
		}
	}
	workspaceId := d.Get("workspace_id").(string)
	version, err := resolveEnforceVersion(d, client, workspaceId)
	if err != nil {
		return err
	}
	createEnforceParams := CreateEnforceRequest{
		Version: version,
	}
	_, err = client.CreateEnforce(createEnforceParams, workspaceId)
	if err != nil {
		return err
	}
	// The workspace is the stable identity for this resource; an enforcement event
	// epoch changes on every enforcement event and cannot be used for reconciliation.
	d.SetId(workspaceId)
	return resourceSecureWorkloadEnforceRead(d, meta)
}

func resourceSecureWorkloadEnforceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	workspaceId := d.Get("workspace_id").(string)
	version, err := resolveEnforceVersion(d, client, workspaceId)
	if err != nil {
		return err
	}
	createEnforceParams := CreateEnforceRequest{
		Version: version,
	}
	// Re-enforce in place; do not disable first.
	_, err = client.CreateEnforce(createEnforceParams, workspaceId)
	if err != nil {
		return err
	}
	return resourceSecureWorkloadEnforceRead(d, meta)
}

func resourceSecureWorkloadEnforceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	workspaceId := d.Id()

	describeApplicationParams := DescribeApplicationRequest{
		ApplicationId: workspaceId,
	}
	application, err := client.DescribeApplication(describeApplicationParams)
	if err != nil {
		if IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if !application.EnforcementEnabled {
		// Enforcement being off means this resource no longer exists conceptually;
		// this lets Terraform recreate it on the next apply.
		d.SetId("")
		return nil
	}

	latestVersion, err := client.LatestPolicyVersion(workspaceId)
	if err != nil {
		return err
	}

	if err := d.Set("workspace_id", workspaceId); err != nil {
		return err
	}
	if err := d.Set("enforcement_enabled", application.EnforcementEnabled); err != nil {
		return err
	}
	enforcedVersion := formatPolicyVersion(application.EnforcedVersion)
	if err := d.Set("enforced_version", enforcedVersion); err != nil {
		return err
	}
	if err := d.Set("latest_version", latestVersion); err != nil {
		return err
	}

	if d.Get("track_latest_version").(bool) {
		// Writing the latest version into `version` is what causes Terraform to plan
		// an in-place update (and thus re-enforce) whenever a newer policy version
		// becomes available on the workspace.
		if err := d.Set("version", latestVersion); err != nil {
			return err
		}
	} else {
		if err := d.Set("version", enforcedVersion); err != nil {
			return err
		}
	}

	return nil
}

func resourceSecureWorkloadEnforceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteEnforce(d.Get("workspace_id").(string))
}
