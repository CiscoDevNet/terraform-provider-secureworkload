// +build all integrationtests

package secureworkload

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	applicationsAPIURL                = os.Getenv("SECUREWORKLOAD_API_URL")
	applicationsAPIKey                = os.Getenv("SECUREWORKLOAD_API_KEY")
	applicationsAPISecret             = os.Getenv("SECUREWORKLOAD_API_SECRET")
	applicationsDefaultAppScopeId     = os.Getenv("SECUREWORKLOAD_APP_SCOPE_ID")
	applicationsDefaultRootScopeAppId = os.Getenv("SECUREWORKLOAD_ROOT_SCOPE_APP_ID")
	applicationsDefaultClientConfig   = Config{
		APIKey:                 applicationsAPIURL,
		APISecret:              applicationsAPIKey,
		APIURL:                 applicationsAPISecret,
		DisableTLSVerification: false,
	}
)

func TestListApplications(t *testing.T) {
	client, err := New(applicationsDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, applicationsDefaultClientConfig)
	}
	_, err = client.ListApplications()
	if err != nil {
		t.Errorf("Error %s listing applications with client %+v", err, client)
	}
}

func TestDeleteApplicationDeletesCreatedApplication(t *testing.T) {
	client, err := New(applicationsDefaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, applicationsDefaultClientConfig)
	}
	createApplicationParams := CreateApplicationRequest{
		AppScopeId:         applicationsDefaultAppScopeId,
		Name:               fmt.Sprintf("test+%d", time.Now().UnixNano()),
		Description:        "TestDeleteApplicationDeletesCreatedApplication",
		AlternateQueryMode: true,
		StrictValidation:   true,
		Primary:            true,
		Clusters: []Cluster{
			Cluster{
				Id:          "ClusterA",
				Name:        "ClusterA",
				Description: "A Cluster.",
				Nodes: []Node{
					Node{
						Name:      "ClusterA Node1",
						IPAddress: "10.0.0.1",
					},
				},
				ConsistentUUID: "ClusterA",
			},
		},
		Filters: []PolicyFilter{
			PolicyFilter{
				Id:   "FilterA",
				Name: "DisplayedClusterName",
				Query: []byte(`{
                      "type": "eq",
                      "field": "ip",
                      "value": "10.0.0.1"
                        }`),
			},
		},
		AbsolutePolicies: []Policy{
			Policy{
				ConsumerFilterId: applicationsDefaultRootScopeAppId,
				ProviderFilterId: applicationsDefaultRootScopeAppId,
				Action:           "ALLOW",
				Layer4NetworkPolicies: []Layer4NetworkPolicy{Layer4NetworkPolicy{
					PortRange: [2]int{80, 80}},
				},
			},
		},
		DefaultPolicies: []Policy{
			Policy{
				ConsumerFilterId: applicationsDefaultRootScopeAppId,
				ProviderFilterId: applicationsDefaultRootScopeAppId,
				Action:           "DENY",
				Layer4NetworkPolicies: []Layer4NetworkPolicy{Layer4NetworkPolicy{
					PortRange: [2]int{8080, 8080}},
				},
			},
		},
		CatchAllAction: "DENY",
	}
	createdApplication, err := client.CreateApplication(createApplicationParams)
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteApplication(createdApplication.Id)
	if err != nil {
		t.Error(err)
	}
	applications, err := client.ListApplications()
	if err != nil {
		t.Errorf("Error %s listing applications with client %+v", err, client)
	}
	for _, application := range applications {
		if application.Id == createdApplication.Id {
			t.Errorf("Application %+v should be deleted but wasn't", application)
		}
	}
}

func TestDescribeApplicationDescribesCreatedApplication(t *testing.T) {
	client, err := New(defaultClientConfig)
	if err != nil {
		t.Errorf("Error %s creating client with config %+v", err, defaultClientConfig)
	}
	createApplicationParams := CreateApplicationRequest{
		AppScopeId:         applicationsDefaultAppScopeId,
		Name:               fmt.Sprintf("test+%d", time.Now().UnixNano()),
		Description:        "TestDescribeApplicationDescribesCreatedApplication",
		AlternateQueryMode: true,
		StrictValidation:   true,
		Primary:            true,
		Clusters: []Cluster{
			Cluster{
				Id:          "ClusterA",
				Name:        "ClusterA",
				Description: "A Cluster.",
				Nodes: []Node{
					Node{
						Name:      "ClusterA Node1",
						IPAddress: "10.0.0.1",
					},
				},
				ConsistentUUID: "ClusterA",
			},
		},
		Filters: []PolicyFilter{
			PolicyFilter{
				Id:   "FilterA",
				Name: "DisplayedClusterName",
				Query: []byte(`{
                      "type": "eq",
                      "field": "ip",
                      "value": "10.0.0.1"
                        }`),
			},
		},
		AbsolutePolicies: []Policy{
			Policy{
				ConsumerFilterId: applicationsDefaultRootScopeAppId,
				ProviderFilterId: applicationsDefaultRootScopeAppId,
				Action:           "ALLOW",
				Layer4NetworkPolicies: []Layer4NetworkPolicy{Layer4NetworkPolicy{
					PortRange: [2]int{80, 80}},
				},
			},
		},
		DefaultPolicies: []Policy{
			Policy{
				ConsumerFilterId: applicationsDefaultRootScopeAppId,
				ProviderFilterId: applicationsDefaultRootScopeAppId,
				Action:           "DENY",
				Layer4NetworkPolicies: []Layer4NetworkPolicy{Layer4NetworkPolicy{
					PortRange: [2]int{8080, 8080}},
				},
			},
		},
		CatchAllAction: "DENY",
	}
	createdApplication, err := client.CreateApplication(createApplicationParams)
	if err != nil {
		t.Error(err)
	}
	describedApplication, err := client.DescribeApplication(DescribeApplicationRequest{
		ApplicationId: createdApplication.Id,
	})
	if describedApplication.Name != createApplicationParams.Name {
		t.Errorf("Expected described application Name to be %s, got %s",
			createApplicationParams.Name, describedApplication.Name)
	}
	if err != nil {
		t.Error(err)
	}
	err = client.DeleteApplication(createdApplication.Id)
	if err != nil {
		t.Error(err)
	}
	applications, err := client.ListApplications()
	if err != nil {
		t.Errorf("Error %s listing applications with client %+v", err, client)
	}
	for _, application := range applications {
		if application.Id == createdApplication.Id {
			t.Errorf("Application %+v should be deleted but wasn't", application)
		}
	}
}
