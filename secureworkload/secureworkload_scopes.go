package secureworkload

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	ScopesAPIV1BasePath = fmt.Sprintf("%s/app_scopes", SecureWorkloadAPIV1BasePath)
)

// ScopeQuery wraps a Filter (or match criteria) associated
// with the scope in conjunction with the filters of
// the parent scopes (all the way to the root scope).
type ScopeQuery struct {
	Type    string       `json:"type"`
	Field   string       `json:"field,omitempty"`
	Value   interface{}  `json:"value,omitempty"`
	Filters []ScopeQuery `json:"filters,omitempty"`
}

// Scope wraps a secureworkload scope attributes including the
// filters used to evaluate the scope and the current state
// of the scope

type GetScope struct {
	Id             string `json:"id"`
	ExactShortName string `json:"short_name,omitempty"`
	ExactName      string `json:"name,omitempty"`
	VRFId          int    `json:"vrf_id,omitempty"`
	RootAppScopeId string `json:"root_app_scope_id,omitempty"`
}
type Scope struct {
	FilterType string `json:"filter_type"`
	// Unique identifier for the scope.
	Id string `json:"id"`
	// Fully qualified name of the scope. This is a fully qualified name; that is, it includes the names of parent scopes (if applicable) all the way to the root scope.
	Name string `json:"name"`
	// User-specified description of the scope.
	Description string `json:"description"`
	// User-specified name for the scope.
	ShortName string `json:"short_name"`
	// Filter (or match criteria) associated with the scope in conjunction with the filters of the parent scopes (all the way to the root scope).
	Query ScopeQuery `json:"query"`
	// Filter (or match criteria) associated with the scope.
	ShortQuery       ScopeQuery `json:"short_query"`
	ParentAppScopeId string     `json:"parent_app_scope_id"`
	RootAppScopeId   string     `json:"root_app_scope_id"`
	// ID of the VRF to which scope belongs.
	VRFId         int    `json:"vrf_id"`
	Priority      string `json:"priority"`
	ShortPriority int    `json:"short_priority"`
	// Used to sort application priorities.
	PolicyPriority int `json:"policy_priority"`
	// Indicates a child or parent query has been updated and that the changes need to be committed.
	Dirty bool `json:"dirty"`
	// Non-null if the query for this scope has been updated but not yet committed.
	DirtyShortQuery  ScopeQuery `json:"dirty_short_query,omitempty"`
	ChildAppScopeIds []string   `json:"child_app_scope_ids"`
	CreatedAt        int64      `json:"created_at"`
	UpdatedAt        int64      `json:"updated_at"`
}

// ShortQuery wraps a query used as part of the request to create a scope
// type ShortQuery struct {
// 	Type  string      `json:"type"`
// 	Field string      `json:"field,omitempty"`
// 	Value interface{} `json:"value,omitempty"`
// }

// CreateScopeRequest wraps parameters for making a request
// to create a scope
type CreateScopeRequest struct {
	// User-specified name for the scope.
	ShortName string `json:"short_name"`

	SubType string `json:"sub_type,omitempty"`
	// User-specified description of the scope.
	Description string `json:"description"`
	// Filter (or match criteria) associated with the scope.
	ShortQuery json.RawMessage `json:"short_query,omitempty"`
	// ID of the parent scope.
	ParentAppScopeId string `json:"parent_app_scope_id"`
	// Used to sort application priorities; default is last.
	PolicyPriority int `json:"policy_priority,omitempty"`
}

func (c Client) GetScopeByParam(getUrl string) ([]GetScope, error) {
	var scope []GetScope
	url := c.Config.APIURL + ScopesAPIV1BasePath + getUrl
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return scope, err
	}
	err = c.Do(request, &scope)
	return scope, err
}

// CreateScope creates a scope with the specified params,
// returning the created scope and error (if any).
func (c Client) CreateScope(params CreateScopeRequest) (Scope, error) {
	var scope Scope
	url := c.Config.APIURL + ScopesAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return scope, err
	}
	err = c.Do(request, &scope)
	return scope, err
}

// DescribeScope describes a scope by id returning the scope
// and error (if any).
func (c Client) DescribeScope(scopeId string) (Scope, error) {
	var scope Scope
	url := c.Config.APIURL + ScopesAPIV1BasePath + fmt.Sprintf("/%s", scopeId)
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return scope, err
	}
	err = c.Do(request, &scope)
	return scope, err
}

// DeleteScope deletes a scope by id returning error (if any).
func (c Client) DeleteScope(scopeId string) error {
	url := c.Config.APIURL + ScopesAPIV1BasePath + fmt.Sprintf("/%s", scopeId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}

// ListScopes lists all scopes readable by the API
// credentials for the given client, returning
// the listed scopes and error (if any)
func (c Client) ListScopes() ([]Scope, error) {
	var scopes []Scope
	url := c.Config.APIURL + ScopesAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return scopes, err
	}
	err = c.Do(request, &scopes)
	return scopes, err
}
