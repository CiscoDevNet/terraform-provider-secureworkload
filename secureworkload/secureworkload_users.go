package secureworkload

import (
	"fmt"
	"net/http"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	UsersAPIV1BasePath = fmt.Sprintf("%s/users", SecureWorkloadAPIV1BasePath)
)

// CreateUserRequest wraps parameters for making a request
// to create a user
type CreateUserRequest struct {
	// Email address associated with the user account.
	Email string `json:"email"`
	// Userʼs first name.
	FirstName string `json:"first_name"`
	// Userʼs last name.
	LastName string `json:"last_name"`
	// (Optional) Root scope to which the user belongs.
	AppScopeId string `json:"app_scope_id,omitempty"`
	// (Optional) A list of roles to be assigned to the user.
	RoleIds []string `json:"role_ids,omitempty"`
}

// User wraps parameters for defining a user along with
// their roles and scope access
type User struct {
	// Unique identifier for the user.
	Id string `json:"id"`
	// Email address associated with the user account.
	Email string `json:"email"`
	// Userʼs first name.
	FirstName string `json:"first_name"`
	// Userʼs last name.
	LastName string `json:"last_name"`
	// (Optional) Root scope to which the user belongs.
	AppScopeId string `json:"app_scope_id,omitempty"`
	// (Optional) A list of roles to be assigned to the user.
	RoleIds []string `json:"role_ids,omitempty"`
	// UNIX timestamp indicating when the user account was disabled. Zero or null if not disabled.
	DisabledAt int `json:"disabled_at,omitempty"`
}

// CreateUser creates a user with the specified params,
// returning the created user and error (if any).
func (c Client) CreateUser(params CreateUserRequest) (User, error) {
	var user User
	url := c.Config.APIURL + UsersAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return user, err
	}
	err = c.Do(request, &user)
	return user, err
}

// DescribeUser describes a user by id returning the user
// and error (if any).
func (c Client) DescribeUser(userId string) (User, error) {
	var user User
	url := c.Config.APIURL + UsersAPIV1BasePath + fmt.Sprintf("/%s", userId)
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return user, err
	}
	err = c.Do(request, &user)
	return user, err
}

// DeleteUser deletes a user by id returning error (if any).
func (c Client) DeleteUser(userId string) error {
	url := c.Config.APIURL + UsersAPIV1BasePath + fmt.Sprintf("/%s", userId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}

// AddRoleToUserRequest wraps the parameters for requesting a role
// be added to a user
type AddRoleToUserRequest struct {
	UserId string `json:"user_id"`
	// ID of the role object to be added
	RoleId string `json:"role_id"`
}

// AddRoleToUser adds the specified role to the user.
// It returns the modified user object and an error, if any
func (c Client) AddRoleToUser(params AddRoleToUserRequest) (User, error) {
	var user User
	url := c.Config.APIURL + UsersAPIV1BasePath + fmt.Sprintf("/%s/add_role", params.UserId)
	request, err := signer.CreateJSONRequest(http.MethodPut, url, params)
	if err != nil {
		return user, err
	}
	err = c.Do(request, &user)
	return user, err
}

// RemoveRoleFromUserRequest wraps the parameters for requesting a role
// be removed from a user
type RemoveRoleFromUserRequest struct {
	// ID of the user
	UserId string `json:"user_id"`
	// ID of the role object to be removed
	RoleId string `json:"role_id"`
}

// RemoveRoleFromUser removes the specified role from the user.
// It returns the modified user object and an error, if any
func (c Client) RemoveRoleFromUser(params RemoveRoleFromUserRequest) (User, error) {
	var user User
	url := c.Config.APIURL + UsersAPIV1BasePath + fmt.Sprintf("/%s/remove_role", params.UserId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, params)
	if err != nil {
		return user, err
	}
	err = c.Do(request, &user)
	return user, err
}

// ListUsersRequest wraps the parameters required to make a
// request to list users
type ListUsersRequest struct {
	// (Optional) Whether to include disabled users; defaults to false.
	IncludeDisabled bool
	// (Optional) Returns only users assigned to the provided scope.
	AppScopeId string
}

// ListUsers lists all users readable by the API
// credentials for the given client, returning
// the listed users and error (if any)
func (c Client) ListUsers(params ListUsersRequest) ([]User, error) {
	var users []User
	url := c.Config.APIURL + UsersAPIV1BasePath
	if params.IncludeDisabled || params.AppScopeId != "" {
		url += "?"
		if params.IncludeDisabled {
			url += "include_disabled=true"
			if params.AppScopeId != "" {
				url += fmt.Sprintf("&app_scope_id=%s", params.AppScopeId)
			}
		}
		if params.AppScopeId != "" {
			url += fmt.Sprintf("app_scope_id=%s", params.AppScopeId)
		}
	}
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return users, err
	}
	err = c.Do(request, &users)
	return users, err
}

// EnableUser makes a request to enable or reactivate a
// deactivated user, returning the user and error (if any).
func (c Client) EnableUser(userId string) (User, error) {
	var user User
	url := c.Config.APIURL + UsersAPIV1BasePath + fmt.Sprintf("/%s/enable", userId)
	request, err := signer.CreateJSONRequest(http.MethodPost, url, nil)
	if err != nil {
		return user, err
	}
	err = c.Do(request, &user)
	return user, err
}
