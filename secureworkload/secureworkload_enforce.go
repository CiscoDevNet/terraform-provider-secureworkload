package secureworkload

import (
	"fmt"
	"net/http"
	"time"

	// "github.com/secureworkload-exchange/terraform-go-sdk/signer"
	"terraform-provider-secureworkload/secureworkload/signer"
)

var (
	EnforceAPIV1BasePath = fmt.Sprintf("%s/applications/", SecureWorkloadAPIV1BasePath)
)

type Enforce struct {
	Id      string `json:"id"`
	Version string `json:"version,omitempty"`
	Epoch   string `json:"epoch,omitempty"`
}

type CreateEnforceRequest struct {
	Version string `json:"version,omitempty"`
	Epoch   string `json:"epoch,omitempty"`
}

func (c Client) CreateEnforce(params CreateEnforceRequest, workspace_id string) (Enforce, error) {
	var enforce Enforce
	url := c.Config.APIURL + EnforceAPIV1BasePath + workspace_id + "/enable_enforce"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return enforce, err
	}
	err = c.Do(request, &enforce)
	return enforce, err
}

func (c Client) DeleteEnforce(workspace_id string) error {
	url := c.Config.APIURL + EnforceAPIV1BasePath + workspace_id + "/disable_enforce"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	c.Do(request, nil)
	time.Sleep(60 * time.Second)
	return nil
}
