// +build all integrationtests

package secureworkload

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	APIURL                  = os.Getenv("SECUREWORKLOAD_API_URL")
	APIKey                  = os.Getenv("SECUREWORKLOAD_API_KEY")
	APISecret               = os.Getenv("SECUREWORKLOAD_API_SECRET")
	defaultParentScopeAppId = os.Getenv("SECUREWORKLOAD_ROOT_SCOPE_APP_ID")
	defaultClientConfig     = Config{
		APIKey:                 APIKey,
		APISecret:              APISecret,
		APIURL:                 APIURL,
		DisableTLSVerification: false,
	}
)

func TestListScopes(t *testing.T) {
	client, err := New(defaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, defaultClientConfig)
	}
	_, err = client.ListScopes()
	if err != nil {
		t.Errorf("Error %s listing scopes with client %+v", err, client)
	}
}

func TestDeleteScopeDeletesCreatedScope(t *testing.T) {
	client, err := New(defaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, defaultClientConfig)
	}
	createScopeParams := CreateScopeRequest{
		ShortName:        fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description:      "TestDeleteScopeDeletesCreatedScope",
		ParentAppScopeId: defaultParentScopeAppId,
		ShortQuery: ShortQuery{
			Type:  "eq",
			Field: "ip",
			Value: "10.0.0.1",
		},
	}
	createdScope, err := client.CreateScope(createScopeParams)
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteScope(createdScope.Id)
	if err != nil {
		t.Error(err)
	}
	scopes, err := client.ListScopes()
	if err != nil {
		t.Errorf("Error %s listing scopes with client %+v", err, client)
	}
	for _, scope := range scopes {
		if scope.Id == createdScope.Id {
			t.Errorf("Scope %+v should be deleted but wasn't", scope)
		}
	}
}

func TestDescribeScopeDescribesCreatedScope(t *testing.T) {
	client, err := New(defaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, defaultClientConfig)
	}
	createScopeParams := CreateScopeRequest{
		ShortName:        fmt.Sprintf("GoSDKTest %d", time.Now().UnixNano()),
		Description:      "TestDescribeScopeDescribesCreatedScope",
		ParentAppScopeId: defaultParentScopeAppId,
		ShortQuery: ShortQuery{
			Type:  "eq",
			Field: "ip",
			Value: "10.0.0.1",
		},
	}
	createdScope, err := client.CreateScope(createScopeParams)
	if err != nil {
		t.Error(err)
	}
	describedScope, err := client.DescribeScope(createdScope.Id)
	if describedScope.ShortName != createScopeParams.ShortName {
		t.Errorf("Expected described scope short name to be %s, got %s",
			createScopeParams.ShortName, describedScope.ShortName)
	}
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteScope(createdScope.Id)
	if err != nil {
		t.Error(err)
	}
	scopes, err := client.ListScopes()
	if err != nil {
		t.Errorf("Error %s listing scopes with client %+v", err, client)
	}
	for _, scope := range scopes {
		if scope.Id == createdScope.Id {
			t.Errorf("Scope %+v should be deleted but wasn't", scope)
		}
	}
}
