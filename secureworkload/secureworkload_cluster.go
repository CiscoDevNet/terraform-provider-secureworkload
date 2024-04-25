package secureworkload

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	ClustersAPIV1BasePath = fmt.Sprintf("%s/applications/", SecureWorkloadAPIV1BasePath)
)

type Clusters struct {
	Id string `json:"id"`

	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	// Query       json.RawMessage `json:"query,omitempty"`
	Approved bool `json:"approved"`

	Query map[string]interface{}

	QueryJSON json.RawMessage `json:"query,omitempty"`
}

type CreateClusterRequest struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Query       json.RawMessage `json:"query,omitempty"`
	Approved    bool            `json:"approved"`
}

func (c Client) CreateCluster(params CreateClusterRequest, workspace_id string) (Clusters, error) {
	var cluster Clusters
	url := c.Config.APIURL + ClustersAPIV1BasePath + workspace_id + "/clusters"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return cluster, err
	}
	err = c.Do(request, &cluster)
	if err != nil {
		return cluster, err
	}
	err = json.Unmarshal(cluster.QueryJSON, &cluster.Query)
	return cluster, err
}

func (c Client) DescribeCluster(clusterId string) (Clusters, error) {
	var cluster Clusters
	url := c.Config.APIURL + SecureWorkloadAPIV1BasePath + "/clusters" + fmt.Sprintf("/%s", clusterId)
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return cluster, err
	}
	err = c.Do(request, &cluster)
	if err != nil {
		return cluster, err
	}
	err = json.Unmarshal(cluster.QueryJSON, &cluster.Query)
	return cluster, err
}

func (c Client) DeleteCluster(workspace_id string, clusterId string) error {
	url := c.Config.APIURL + SecureWorkloadAPIV1BasePath + "/clusters" + fmt.Sprintf("/%s", clusterId)
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}

func (c Client) ListCluster(workspace_id string) ([]Clusters, error) {
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
