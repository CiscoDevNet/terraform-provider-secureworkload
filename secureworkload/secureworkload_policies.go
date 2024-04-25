package secureworkload

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	PolicyAPIV1BasePath = fmt.Sprintf("%s/applications/", SecureWorkloadAPIV1BasePath)
)

type Policies struct {
	Id         string `json:"id"`
	ConsumerId string `json:"consumer_filter_id"`
	ProviderId string `json:"provider_filter_id"`
	Version    string `json:"version,omitempty"`
	Rank       string `json:"rank,omitempty"`
	Action     string `json:"policy_action"`
	Priority   int    `json:"priority,omitempty"`
}

type CreatePolicyRequest struct {
	ConsumerId string `json:"consumer_filter_id"`
	ProviderId string `json:"provider_filter_id"`
	Version    string `json:"version,omitempty"`
	Rank       string `json:"rank,omitempty"`
	Action     string `json:"policy_action"`
	Priority   int    `json:"priority,omitempty"`
	// Layer4NetworkPolicies []Layer4Network `json:"l4_params"`
}

// type Layer4Network struct {
// 	Protocol  int    `json:"proto,omitempty"`
// 	PortRange [2]int `json:"port"`
// 	Approved  bool   `json:"approved"`
// }

func (c Client) CreatePolicy(params CreatePolicyRequest, workspace_id string) (Policies, error) {
	var policy Policies
	url := c.Config.APIURL + PolicyAPIV1BasePath + workspace_id + "/policies"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return policy, err
	}
	err = c.Do(request, &policy)
	return policy, err
}

func (c Client) DescribePolicy(policyId string) (Policies, error) {
	var policy Policies
	url := c.Config.APIURL + SecureWorkloadAPIV1BasePath + "/policies" + fmt.Sprintf("/%s", policyId)
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return policy, err
	}
	err = c.Do(request, &policy)
	return policy, err
}

func (c Client) DeletePolicy(workspace_id string, policyId string) error {
	url := c.Config.APIURL + SecureWorkloadAPIV1BasePath + "/policies" + fmt.Sprintf("/%s", policyId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}

func (c Client) ListPolicy(workspace_id string) ([]Clusters, error) {
	var clusters []Clusters
	url := c.Config.APIURL + ClustersAPIV1BasePath + workspace_id + "/clusters"
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return clusters, err
	}
	err = c.Do(request, &clusters)
	if err != nil {
		return clusters, err
	}
	for _, filter := range clusters {
		err = json.Unmarshal(filter.QueryJSON, &filter.Query)
		if err != nil {
			return clusters, err
		}
	}
	return clusters, err
}
