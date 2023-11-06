// +build all integrationtests

package secureworkload

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	filtersAPIURL                  = os.Getenv("SECUREWORKLOAD_API_URL")
	filtersAPIKey                  = os.Getenv("SECUREWORKLOAD_API_KEY")
	filtersAPISecret               = os.Getenv("SECUREWORKLOAD_API_SECRET")
	filtersDefaultParentScopeAppId = os.Getenv("SECUREWORKLOAD_ROOT_SCOPE_APP_ID")
	filtersDefaultClientConfig     = Config{
		APIKey:                 filtersAPIKey,
		APISecret:              filtersAPISecret,
		APIURL:                 filtersAPIURL,
		DisableTLSVerification: false,
	}
)

func TestListFilters(t *testing.T) {
	client, err := New(filtersDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, filtersDefaultClientConfig)
	}
	_, err = client.ListFilters()
	if err != nil {
		t.Errorf("Error %s listing filters with client %+v", err, client)
	}
}

func TestDeleteFilterDeletesCreatedFilter(t *testing.T) {
	client, err := New(filtersDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, filtersDefaultClientConfig)
	}
	createFilterParams := CreateFilterRequest{
		Name:       fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		AppScopeId: filtersDefaultParentScopeAppId,
		Query: []byte(`{
	                  "type": "eq",
	                  "field": "ip",
	                  "value": "10.0.0.1"
	                    }`),
	}
	createdFilter, err := client.CreateFilter(createFilterParams)
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteFilter(createdFilter.Id)
	if err != nil {
		t.Error(err)
	}
	filters, err := client.ListFilters()
	if err != nil {
		t.Errorf("Error %s listing filters with client %+v", err, client)
	}
	for _, filters := range filters {
		if filters.Id == createdFilter.Id {
			t.Errorf("Filters %+v should be deleted but wasn't", filters)
		}
	}
}

func TestDescribeFilterDescribesCreatedFilter(t *testing.T) {
	client, err := New(filtersDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, filtersDefaultClientConfig)
	}
	createFilterParams := CreateFilterRequest{
		Name:       fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		AppScopeId: filtersDefaultParentScopeAppId,
		Query: []byte(`{
	                  "type": "eq",
	                  "field": "ip",
	                  "value": "10.0.0.1"
	                    }`),
	}
	createdFilter, err := client.CreateFilter(createFilterParams)
	if err != nil {
		t.Error(err)
	}
	describedFilter, err := client.DescribeFilter(createdFilter.Id)
	if describedFilter.Name != createFilterParams.Name {
		t.Errorf("Expected described filter short name to be %s, got %s",
			createFilterParams.Name, describedFilter.Name)
	}
	err = client.DeleteFilter(createdFilter.Id)
	if err != nil {
		t.Error(err)
	}
	filters, err := client.ListFilters()
	if err != nil {
		t.Errorf("Error %s listing filters with client %+v", err, client)
	}
	for _, filters := range filters {
		if filters.Id == createdFilter.Id {
			t.Errorf("Filters %+v should be deleted but wasn't", filters)
		}
	}
}

func TestCreateFilterAllowsForCreatingFiltersWithComplexQuery(t *testing.T) {
	client, err := New(filtersDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, filtersDefaultClientConfig)
	}
	createFilterParams := CreateFilterRequest{
		Name:       fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		AppScopeId: filtersDefaultParentScopeAppId,
		Query: []byte(`{
               "type": "and",
               "filters": [
                  {
                     "field": "vrf_id",
                     "type": "eq",
                     "value": 700056
                  },
                  {
                     "type": "or",
                     "filters": [
                        {
                           "field": "ip",
                           "type": "eq",
                           "value": "10.254.252.43"
                        },
                        {
                           "field": "ip",
                           "type": "eq",
                           "value": "10.254.252.51"
                        },
                        {
                           "field": "ip",
                           "type": "eq",
                           "value": "10.254.252.52"
                        }
                     ]
                  }
               ]
            }`),
	}
	createdFilter, err := client.CreateFilter(createFilterParams)
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteFilter(createdFilter.Id)
	if err != nil {
		t.Error(err)
	}
}
