package secureworkload

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	// client "github.com/secureworkload-exchange/terraform-go-sdk"
	// secureworkload "github.com/secureworkload-exchange/terraform-go-sdk"
)

var (
	// ValidAbilities        = []string{"SCOPE_READ", "SCOPE_WRITE", "EXECUTE", "ENFORCE", "SCOPE_OWNER", "DEVELOPER"}
	// AccessTypeDescription = fmt.Sprintf("The type of access to grant the role to the `access_app_scope_id` scope.\n Valid values are [%s]", strings.Join(ValidAbilities, ", "))
)

func resourceSecureWorkloadRole() *schema.Resource {
	return &schema.Resource{
		Description: "Resource for creating a role in Secure Workload\n" +
			"\n" +
			"## Example\n" +
			"An example is shown below: \n" +
			"```hcl\n" +
			"resource \"secureworkload_role\" \"role1\" {\n" +
			"	 app_scope_id = data.secureworkload_scope.scope.id\n" +
			"	 access_app_scope_id = data.secureworkload_scope2.scope2.id\n" +
			"    name = \"read_role\"\n" +
			"    access_type = \"scope_read\"\n" +
			"    user_ids = [\"<USER_IDS>\"]\n" +
			"    description = \"Demo description for role\"\n" +
			"}\n" +
			"```\n" +
			"**Note:** If creating multiple rules during a single `terraform apply`, remember to use `depends_on` to chain the rules so that terraform creates it in the same order that you intended.\n" ,
		Create: resourceSecureWorkloadRoleCreate,
		Update: nil,
		Read:   resourceSecureWorkloadRoleRead,
		Delete: resourceSecureWorkloadRoleDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "(Optional) User-specified name for the role.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The role's description",
			},
			"access_app_scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The scope to which this role will be given access",
			},
			"app_scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The scope in which this role will be created",
			},
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  `The type of access to grant the role to the "access_app_scope_id" scope.\n Valid values are SCOPE_READ", "SCOPE_WRITE", "EXECUTE", "ENFORCE", "SCOPE_OWNER", "DEVELOPER"`,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := strings.ToUpper(val.(string))
					allowedValues := []string{"SCOPE_READ", "SCOPE_WRITE", "EXECUTE", "ENFORCE", "SCOPE_OWNER", "DEVELOPER"}
					for _, allowed := range allowedValues {
						if v == allowed {
							return
						}
					}
					errs = append(errs, fmt.Errorf("%q must be in %v, got: %q", key, allowedValues, v))
					return
				},
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The users to which this role will be assigned",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSecureWorkloadRoleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(Client)
	tfUserIds := d.Get("user_ids").(*schema.Set).List()
	userIds := []string{}
	for _, tfUserId := range tfUserIds {
		if tfUserId != "" {
			userIds = append(userIds, tfUserId.(string))
		}
	}

	createScopedRoleForUsersParams := CreateScopedRoleForUsersRequest{
		CreateScopedRoleRequest: CreateScopedRoleRequest{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			AppScopeId:          d.Get("app_scope_id").(string),
			AbilitiesAppScopeId: d.Get("access_app_scope_id").(string),
			Ability:             d.Get("access_type").(string),
		},
		Users: userIds,
	}

	response, err := client.CreateScopedRoleForUsers(createScopedRoleForUsersParams)
	if err != nil {
		return err
	}
	d.SetId(response.RoleId)
	return nil
}

func resourceSecureWorkloadRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	role, err := client.GetRole(d.Id())
	if err != nil {
		return err
	}
	d.Set("app_scope_id", role.AppScopeId)
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	return nil
}
func resourceSecureWorkloadRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteRole(d.Id())
}
