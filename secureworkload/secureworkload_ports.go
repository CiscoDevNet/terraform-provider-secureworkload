package secureworkload

import (
	"fmt"
	"net/http"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	PortsAPIV1BasePath = fmt.Sprintf("%s/policies/", SecureWorkloadAPIV1BasePath)
)

type Port struct {
	Id          string `json:"id"`
	StartPort   int    `json:"start_port"`
	EndPort     int    `json:"end_port"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Proto       int    `json:"proto,omitempty"`
}

type CreatePortRequest struct {
	StartPort   int    `json:"start_port"`
	EndPort     int    `json:"end_port"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Proto       int    `json:"proto,omitempty"`
}

func (c Client) CreatePort(params CreatePortRequest, policy_id string) (Port, error) {
	var port Port
	url := c.Config.APIURL + PortsAPIV1BasePath + policy_id + "/l4_params"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return port, err
	}
	err = c.Do(request, &port)
	return port, err
}

func (c Client) DescribePort(policy_id string, portId string) (Port, error) {
	var port Port
	url := c.Config.APIURL + SecureWorkloadAPIV1BasePath + "/l4_params" + fmt.Sprintf("/%s", portId)
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return port, err
	}
	err = c.Do(request, &port)
	return port, err
}

func (c Client) DeletePort(policy_id string, portId string) error {
	url := c.Config.APIURL + PortsAPIV1BasePath + policy_id + "/l4_params" + fmt.Sprintf("/%s", portId)
	_, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return nil
}
