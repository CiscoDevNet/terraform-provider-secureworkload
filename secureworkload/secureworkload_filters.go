package secureworkload

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	FiltersAPIV1BasePath = fmt.Sprintf("%s/filters/inventories", SecureWorkloadAPIV1BasePath)
)

// Filter wraps a secureworkload inventory filter.
// Inventory filters encode the match criteria for inventory-search queries,
// used to determine the scope based policies to apply to an application.
type Filter struct {
	// Unique identifier for the filter.
	Id string `json:"id"`
	// User-specified name for the inventory filter.
	Name string `json:"name"`
	// ID of the scope associated with the filter.
	AppScopeId string `json:"app_scope_id"`
	// Filter (or match criteria) associated with the filter.
	ShortQuery ScopeQuery `json:"short_query"`
	// When true, the filter is restricted to the ownership scope.
	Primary bool `json:"primary"`
	// When true the filter provides a service for its scope. Must also be primary/scope restricted.
	Public bool `json:"public"`
	// Filter (or match criteria) associated with the filter in conjunction with the filters of the parent scopes. These conjunctions take effect if ‘restricted to ownership scope’ checkbox is checked. If primary field is false then query is same as ShortQuery.
	Query map[string]interface{}
	// Raw query returned by the SecureWorkload API over the wire.
	QueryJSON json.RawMessage `json:"query,omitempty"`
}

// CreateFilterRequest wraps parameters for making a request to create a filter.
type CreateFilterRequest struct {
	// User-specified name for the inventory filter.
	Name string `json:"name"`
	// ID of the scope associated with the filter.
	AppScopeId string `json:"app_scope_id"`
	// Filter (or match criteria) associated with the scope.
	Query json.RawMessage `json:"query,omitempty"`
	// When true, the filter is restricted to the ownership scope.
	Primary bool `json:"primary"`
	// When true the filter provides a service for its scope. Must also be primary/scope restricted.
	Public bool `json:"public"`
}

func (c Client) GetFilterByParam(getUrl string) ([]Application, error) {
	var filter []Application
	url := c.Config.APIURL + FiltersAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return filter, err
	}
	err = c.Do(request, &filter)
	return filter, err
}

// CreateFilter creates a filters with the specified params,
// returning the created filters and error (if any).
func (c Client) CreateFilter(params CreateFilterRequest) (Filter, error) {
	var filter Filter
	url := c.Config.APIURL + FiltersAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return filter, err
	}
	err = c.Do(request, &filter)
	if err != nil {
		return filter, err
	}
	err = json.Unmarshal(filter.QueryJSON, &filter.Query)
	return filter, err
}

// DescribeFilter describes a filter by id
// returning the filter and error (if any).
func (c Client) DescribeFilter(filterId string) (Filter, error) {
	var filter Filter
	url := c.Config.APIURL + FiltersAPIV1BasePath + fmt.Sprintf("/%s", filterId)
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return filter, err
	}
	err = c.Do(request, &filter)
	if err != nil {
		return filter, err
	}
	err = json.Unmarshal(filter.QueryJSON, &filter.Query)
	return filter, err
}

// DeleteFilter deletes a filter by id returning error (if any).
func (c Client) DeleteFilter(filterId string) error {
	url := c.Config.APIURL + FiltersAPIV1BasePath + fmt.Sprintf("/%s", filterId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}

// ListFilters lists all filters readable by the API
// credentials for the given client, returning
// the listed filters and error (if any)
func (c Client) ListFilters() ([]Filter, error) {
	var filters []Filter
	url := c.Config.APIURL + FiltersAPIV1BasePath
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return filters, err
	}
	err = c.Do(request, &filters)
	if err != nil {
		return filters, err
	}
	for _, filter := range filters {
		err = json.Unmarshal(filter.QueryJSON, &filter.Query)
		if err != nil {
			return filters, err
		}
	}
	return filters, err
}
