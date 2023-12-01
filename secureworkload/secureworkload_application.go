package secureworkload

import (
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-secureworkload/secureworkload/signer"
	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
)

var (
	ApplicationsAPIV1BasePath = fmt.Sprintf("%s/applications", SecureWorkloadAPIV1BasePath)
)

// Application wraps parameters for defining, analyzing and enforcing policies for a
// particular application.
type Application struct {
	// Unique identifier for the application.
	Id string `json:"id"`
	// User-specified name for the application.
	Name string `json:"name"`
	// User-specified description of the application.
	Description string `json:"description"`
	// ID of the scope assigned to the application.
	AppScopeId string `json:"app_scope_id"`
	// First and last name of the user who created the application.
	Author string `json:"author"`
	// Indicates if the application is primary for its scope.
	Primary bool `json:"primary"`
	// Indicates if “dynamic mode” is used for the application. In dynamic mode, an ADM run creates one or more candidate queries for each cluster. Default value is true.
	AlternateQueryMode bool `json:"alternate_query_mode"`
	// Unix timestamp indicating when the application was created.
	CreatedAt int `json:"created_at"`
	// The latest adm (v*) version of the application.
	LatestADMVersion int `json:"latest_adm_version"`
	// Indicates if enforcement is enabled on the application.
	EnforcementEnabled bool `json:"enforcement_enabled"`
	// The enforced p* version of the application.
	EnforcedVersion int `json:"enforced_version"`
}

// CreateApplicationRequest wraps parameters for making a request to create a application.
type CreateApplicationRequest struct {
	// ID of the scope assigned to the application.
	AppScopeId string `json:"app_scope_id"`
	// (Optional) User-specified name for the application.
	Name string `json:"name,omitempty"`
	// (Optional) User-specified description of the application.
	Description string `json:"description,omitempty"`
	// (Optional) Indicates if “dynamic mode” is used for the application. In dynamic mode, an ADM run creates one or
	// more candidate queries for each cluster. Default value is true.
	AlternateQueryMode bool `json:"alternate_query_mode,omitempty"`
	// (Optional) Return an error if there are unknown keys/attributes in the uploaded data. Useful for catching misspelled keys. Default value is false.
	StrictValidation bool `json:"strict_validation,omitempty"`
	// (Optional) Set to true to indicate this application is primary for the given scope. Default value is true.
	Primary bool `json:"primary"`
	// Groups of nodes to be used to define policies.
	Clusters []Cluster `json:"clusters"`
	// Filters on data center assets.
	Filters []PolicyFilter `json:"inventory_filters"`
	// Ordered policies to be created with the absolute rank.
	AbsolutePolicies []Policy `json:"absolute_policies"`
	// Ordered policies to be created with the default rank.
	DefaultPolicies []Policy `json:"default_policies"`
	// “ALLOW” or “DENY”
	CatchAllAction string `json:"catch_all_action"`
}

// PolicyFilter wrap a collection of
// inventory filters on data center assets.
// used to define an application policy.
type PolicyFilter struct {
	// Unique identifier to be used with policies.
	Id string `json:"id"`
	// Displayed name of the cluster.
	Name string `json:"name"`
	// JSON object representation of an inventory filter query.
	Query json.RawMessage `json:"query"`
}

// Cluster wraps a groups of nodes to be used to define policies.
type Cluster struct {
	//     Unique identifier to be used with policies.
	Id string `json:"id"`
	// Cluster display name.
	Name string `json:"name"`
	// Description of the cluster.
	Description string `json:"description"`
	// Nodes or endpoints that are part of the cluster.
	Nodes []Node `json:"nodes"`
	// Must be unique to a given application. After an ADM run, the similar/same clusters in the next version will maintain the consistent_uuid.
	ConsistentUUID string `json:"consistent_uuid"`
}

// Node represents an endpoint that is part of a cluster
type Node struct {
	// IP address or subnet of the node; for example, 10.0.0.1/8 or 1.2.3.4.
	IPAddress string `json:"ip"`
	// Displayed name of the node.
	Name string `json:"name"`
}

// Policy describes an application level policy, either absolute or default.
type Policy struct {
	// ID of a cluster, user inventory filter, or application scope.
	ConsumerFilterId string `json:"consumer_filter_id"`
	// ID of a cluster, user inventory filter, or application scope.
	ProviderFilterId string `json:"provider_filter_id"`
	// “ALLOW” or “DENY”
	Action string `json:"action"`
	// List of allowed ports and protocols.
	Layer4NetworkPolicies []Layer4NetworkPolicy `json:"l4_params"`
}

// Layer4NetworkPolicy wraps parameters for
// enforcing a layer 4 networking policy based off a flows protocol and ports.
type Layer4NetworkPolicy struct {
	// Protocol integer value (NULL means all protocols).
	Protocol int `json:"proto,omitempty"`
	// Inclusive range of ports; for example, [80, 80] or [5000, 6000].
	PortRange [2]int `json:"port"`
	// (Optional) Indicates whether the policy is approved. Default is false.
	Approved bool `json:"approved"`
}
func (c Client) GetApplicationByParam(getUrl string) ([]Application, error) {
	var scope []Application
	url := c.Config.APIURL + ApplicationsAPIV1BasePath + getUrl
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return scope, err
	}
	err = c.Do(request, &scope)
	return scope, err
}
// CreateApplication creates a application with
// the specified params, returning the created application and error (if any).
func (c Client) CreateApplication(params CreateApplicationRequest) (Application, error) {
	var application Application
	url := c.Config.APIURL + ApplicationsAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return application, err
	}
	err = c.Do(request, &application)
	return application, err
}

// DescribeApplicationRequest wraps parameters for making a request to describe a application.
type DescribeApplicationRequest struct {
	ApplicationId string
	// (Optional) A version ID in the form “v10” or “p10”; defaults to “latest.”
	Version string
}

// DescribeApplication describes a application
// by id and version (defaulting to latest)
// returning the application and error (if any).
func (c Client) DescribeApplication(params DescribeApplicationRequest) (Application, error) {
	var application Application
	url := c.Config.APIURL + ApplicationsAPIV1BasePath + fmt.Sprintf("/%s", params.ApplicationId)
	if params.Version != "" {
		url += fmt.Sprintf("/versions/%s", params.Version)
	}
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return application, err
	}
	err = c.Do(request, &application)
	return application, err
}

// DeleteApplication deletes a application by id returning error (if any).
func (c Client) DeleteApplication(applicationId string) error {
	url := c.Config.APIURL + ApplicationsAPIV1BasePath + fmt.Sprintf("/%s", applicationId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}

// ListApplications lists all applications readable by the API
// credentials for the given client, returning
// the listed applications and error (if any)
func (c Client) ListApplications(app_scope_id string) ([]Application, error) {
	var applications []Application
	appendUrl := "?app_scope_id=" + app_scope_id
	url := c.Config.APIURL + ApplicationsAPIV1BasePath + appendUrl
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return applications, err
	}
	err = c.Do(request, &applications)
	for _, application := range applications {
		if application.Primary {
			return applications, err
		}
	}
	return nil , err
}
