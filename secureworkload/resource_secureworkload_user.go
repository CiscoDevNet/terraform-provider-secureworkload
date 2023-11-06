package secureworkload

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// client "github.com/secureworkload-exchange/terraform-go-sdk"
	// secureworkload "github.com/secureworkload-exchange/terraform-go-sdk"
)

func resourceSecureWorkloadUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecureWorkloadUserCreate,
		Update: nil,
		Read:   resourceSecureWorkloadUserRead,
		Delete: resourceSecureWorkloadUserDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Email address associated with the user account.",
			},
			"first_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Userʼs first name.",
			},
			"last_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Userʼs last name.",
			},
			"app_scope_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "(Optional) Root scope to which the user belongs.",
			},
			"role_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "(Optional) A list of roles to be assigned to the user.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_existing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true, and an existing but disabled user with the same email exists they will be enabled.",
			},
			"disabled_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "UNIX timestamp indicating when the user account was disabled. Zero or null if not disabled.",
			},
		},
	}
}

func resourceSecureWorkloadUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	enableExistingUser := d.Get("enable_existing").(bool)
	if enableExistingUser {
		users, err := client.ListUsers(ListUsersRequest{
			AppScopeId:      d.Get("app_scope_id").(string),
			IncludeDisabled: true,
		})
		if err != nil {
			return err
		}
		var userExists bool
		var user User
		for _, existingUser := range users {
			if existingUser.Email == d.Get("email").(string) {
				user = existingUser
				userExists = true
				break
			}
		}
		if userExists {
			user, err = client.EnableUser(user.Id)
			if err != nil {
				return err
			}
			d.SetId(user.Id)
			return nil
		}
	}
	tfRoleIds := d.Get("role_ids").(*schema.Set).List()
	roleIds := make([]string, len(tfRoleIds))
	for _, tfRoleId := range tfRoleIds {
		roleIds = append(roleIds, tfRoleId.(string))
	}
	createUserParams := CreateUserRequest{
		Email:      d.Get("email").(string),
		FirstName:  d.Get("first_name").(string),
		LastName:   d.Get("last_name").(string),
		AppScopeId: d.Get("app_scope_id").(string),
		RoleIds:    roleIds,
	}
	user, err := client.CreateUser(createUserParams)
	if err != nil {
		return err
	}
	d.SetId(user.Id)
	return nil
}

func resourceSecureWorkloadUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	user, err := client.DescribeUser(d.Id())
	if err != nil {
		return err
	}
	d.Set("email", user.Email)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("app_scope_id", user.AppScopeId)
	d.Set("role_ids", user.RoleIds)
	d.Set("disabled_at", user.DisabledAt)
	return nil
}

func resourceSecureWorkloadUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(Client)
	return client.DeleteUser(d.Id())
}
