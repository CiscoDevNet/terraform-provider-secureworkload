package secureworkload

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	RolesAPIV1BasePath   = fmt.Sprintf("%s/roles", SecureWorkloadAPIV1BasePath)
	ValidAbilities2       = []string{"SCOPE_READ", "SCOPE_WRITE", "EXECUTE", "ENFORCE", "SCOPE_OWNER", "DEVELOPER"}
	ValidAbilitiesString = strings.Join(ValidAbilities2, ", ")
)

func isAbilityValid(toValidate string) bool {
	for _, ability := range ValidAbilities2 {
		if ability == toValidate {
			return true
		}
	}
	return false
}

// Role wraps the parameters required to define a role and grant it
// scope access abilities
type Role struct {
	// Unique identifier for the role.
	Id string `json:"id"`
	// (Optional) Application for which the scope is defined
	AppScopeId string `json:"app_scope_id,omitempty"`
	// User-specified name for the role
	Name string `json:"name"`
	// User-specified description for the role
	Description string `json:"description"`
}

func (c Client) GetRoleByParam(getUrl string) ([]Role, error) {
	var role []Role
	url := c.Config.APIURL + RolesAPIV1BasePath + getUrl
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return role, err
	}
	err = c.Do(request, &role)
	return role, err
}

// ListRoles lists all roles readable by the API
// credentials for the given client, returning
// the listed roles and error (if any)
func (c Client) ListRoles() ([]Role, error) {
	var roles []Role
	url := c.Config.APIURL + RolesAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return roles, err
	}
	err = c.Do(request, &roles)
	return roles, err
}

// DeleteRole deletes a role by id, returning error (if any).
func (c Client) DeleteRole(roleId string) error {
	url := fmt.Sprintf("%s%s/%s", c.Config.APIURL, RolesAPIV1BasePath, roleId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}

// CreateRoleRequest wraps parameters for making a request
// to create a role
type CreateRoleRequest struct {
	// User-specified name for the role
	Name string `json:"name"`
	// User-specified description for the role
	Description string `json:"description"`
	// (Optional) Application for which the scope is defined
	AppScopeId string `json:"app_scope_id,omitempty"`
}

// CreateRole creates a role with the specified params,
// returning the created role and error (if any).
func (c Client) CreateRole(params CreateRoleRequest) (Role, error) {
	var role Role
	url := c.Config.APIURL + RolesAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return role, err
	}
	err = c.Do(request, &role)
	return role, err
}

// GetRole describes a role by id returning the role
// and error (if any).
func (c Client) GetRole(roleId string) (Role, error) {
	var role Role
	url := fmt.Sprintf("%s%s/%s", c.Config.APIURL, RolesAPIV1BasePath, roleId)

	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return role, err
	}

	err = c.Do(request, &role)
	return role, err
}

// GiveScopeAccessToRoleRequest wraps the parameters to
// make a request to give scope access to a role
type GiveScopeAccessToRoleRequest struct {
	// The role to which access will be given
	RoleId string `json:"role_id"`
	// The app scope to which the role will be given access
	AppScopeId string `json:"app_scope_id"`
	// Possible values are SCOPE_READ, SCOPE_WRITE, EXECUTE, ENFORCE, SCOPE_OWNER, DEVELOPER
	Ability string `json:"ability"`
}

// RoleScopeResponse wraps the response received from
// adding scope access to a role
type RoleScopeResponse struct {
	// The added scope id
	AppScopeId string `json:"app_scope_id"`
	// The role to which the access was added
	RoleId string `json:"role_id"`
	// The type of access ability granted
	Ability string `json:"ability"`
	// Indicates whether or not the access was inherited
	Inherited bool `json:"inherited"`
}

// GiveScopeAccessToRole gives a role a specific level of access to a scope.
// It returns a RoleScopeResponse and an error, if any
func (c Client) GiveScopeAccessToRole(params GiveScopeAccessToRoleRequest) (RoleScopeResponse, error) {
	var roleScopeResponse RoleScopeResponse
	uppercasedAbility := strings.ToUpper(params.Ability)
	if !isAbilityValid(uppercasedAbility) {
		return roleScopeResponse, errors.New(fmt.Sprintf("Ability %s not valid. Valid abilities are %s", uppercasedAbility, ValidAbilitiesString))
	}
	params.Ability = uppercasedAbility

	url := fmt.Sprintf("%s%s/%s/capabilities", c.Config.APIURL, RolesAPIV1BasePath, params.RoleId)
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return roleScopeResponse, err
	}

	err = c.Do(request, &roleScopeResponse)
	if err != nil {
		return roleScopeResponse, err
	}
	return roleScopeResponse, err
}

// CreateScopedRoleRequest wraps the parameters to make a request
// to create a role with a specified app scope access
type CreateScopedRoleRequest struct {
	// User-specified name for the role
	Name string `json:"name"`
	// User-specified description for the role
	Description string `json:"description"`
	// The application scope under which the role will be created
	AppScopeId string `json:"app_scope_id"`
	// The application scope which will be applied to the role
	AbilitiesAppScopeId string `json:"abilities_app_scope_id"`
	// Possible values are SCOPE_READ, SCOPE_WRITE, EXECUTE, ENFORCE, SCOPE_OWNER, DEVELOPER
	Ability string `json:"ability"`
}

type CreateScopedRoleResponse struct {
	RoleScopeResponse
	// The role which was created
	Role Role `json:"role"`
}

// CreateScopedRole creates a role and gives it the specified scope access.
// It returns the created role and an error, if any.
func (c Client) CreateScopedRole(params CreateScopedRoleRequest) (CreateScopedRoleResponse, error) {

	var role Role

	role, err := c.CreateRole(CreateRoleRequest{
		Name:        params.Name,
		Description: params.Description,
		AppScopeId:  params.AppScopeId,
	})
	fullResponse := CreateScopedRoleResponse{Role: role}
	if err != nil {
		return fullResponse, err
	}

	giveScopeAccessParams := GiveScopeAccessToRoleRequest{
		RoleId:     role.Id,
		AppScopeId: params.AbilitiesAppScopeId,
		Ability:    params.Ability,
	}

	roleScopeResponse, err := c.GiveScopeAccessToRole(giveScopeAccessParams)
	if err != nil {
		defer c.DeleteRole(role.Id)
		return fullResponse, err
	}
	fullResponse.RoleScopeResponse = roleScopeResponse
	return fullResponse, err
}

// CreateScopedRoleForUsersRequest wraps the parameters to make a request
// to create a role with a specified app scope access which is assigned to users
type CreateScopedRoleForUsersRequest struct {
	CreateScopedRoleRequest
	// Users to which the created role will be assigned
	Users []string `json:"users"`
}

type CreateScopedRoleForUsersResponse struct {
	CreateScopedRoleResponse
	// The modified user objects reflecting the role addition
	Users []User `json:"users"`
}

// CreateScopedRoleForUsers creates a role, gives it the specified scope access,
// then assigns it to users.
// It returns the created role and an error, if any.
func (c Client) CreateScopedRoleForUsers(params CreateScopedRoleForUsersRequest) (CreateScopedRoleForUsersResponse, error) {
	createScopedRoleResponse, err := c.CreateScopedRole(params.CreateScopedRoleRequest)

	fullResponse := CreateScopedRoleForUsersResponse{CreateScopedRoleResponse: createScopedRoleResponse}
	if err != nil {
		return fullResponse, err
	}

	for _, userId := range params.Users {
		user, err := c.AddRoleToUser(AddRoleToUserRequest{
			UserId: userId,
			RoleId: fullResponse.Role.Id,
		})
		// TODO: what about a map of userIds to an object with either the response or resulting error for that user
		// it's a heavier lift but provides more info about individual failures
		fullResponse.Users = append(fullResponse.Users, user)
		if err != nil {
			defer c.DeleteRole(fullResponse.Role.Id)
			return fullResponse, err
		}
	}

	return fullResponse, err
}
