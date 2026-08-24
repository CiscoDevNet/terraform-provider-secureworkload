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
	Id        string `json:"id"`
	StartPort int    `json:"start_port"`
	EndPort   int    `json:"end_port"`
	// PortRange captures the shape the API actually returns for an
	// l4_param: "port": [start, end]. The API does NOT echo back
	// start_port/end_port on reads, even though it accepts them on
	// create, so anything that inspects a returned l4_param must go
	// through normalize() before comparing ports.
	PortRange   []int  `json:"port,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Proto       int    `json:"proto,omitempty"`
}

// normalize backfills StartPort/EndPort from the API's "port" array so the
// rest of the provider can treat a decoded l4_param uniformly.
func (p *Port) normalize() {
	if len(p.PortRange) > 0 && p.StartPort == 0 {
		p.StartPort = p.PortRange[0]
	}
	if len(p.PortRange) > 1 && p.EndPort == 0 {
		p.EndPort = p.PortRange[1]
	}
}

type CreatePortRequest struct {
	StartPort   int    `json:"start_port"`
	EndPort     int    `json:"end_port"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Proto       int    `json:"proto,omitempty"`
}

// policyWithL4Params is a decoding-only shape used to capture the parent
// policy object that the API returns from the create/get l4_param
// endpoints. The API responds with the updated parent policy (top-level
// "id" is the POLICY id, not the l4_param id), so we need access to the
// "l4_params" list to find the actual port entry we care about.
type policyWithL4Params struct {
	Id       string `json:"id"`
	L4Params []Port `json:"l4_params"`
}

// findMatchingL4Param looks for the last (newest) entry in l4Params that
// matches the requested proto/start_port/end_port. Returns the matching
// Port and true if found.
func findMatchingL4Param(l4Params []Port, params CreatePortRequest) (Port, bool) {
	var match Port
	found := false
	for i := range l4Params {
		p := l4Params[i]
		p.normalize()
		if p.StartPort == params.StartPort && p.EndPort == params.EndPort && p.Proto == params.Proto {
			match = p
			found = true
		}
	}
	return match, found
}

func (c Client) CreatePort(params CreatePortRequest, policy_id string) (Port, error) {
	var port Port
	var resp policyWithL4Params
	url := c.Config.APIURL + PortsAPIV1BasePath + policy_id + "/l4_params"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return port, err
	}
	err = c.Do(request, &resp)
	if err != nil {
		return port, err
	}

	if len(resp.L4Params) > 0 {
		if match, found := findMatchingL4Param(resp.L4Params, params); found {
			return match, nil
		}
	}

	// Fall back to fetching the parent policy directly, in case the
	// create response did not include the l4_params list.
	policyResp, err := c.DescribePolicyL4Params(policy_id)
	if err != nil {
		return port, err
	}
	if match, found := findMatchingL4Param(policyResp.L4Params, params); found {
		return match, nil
	}

	return port, fmt.Errorf("could not determine the l4_param id for the newly created port (proto=%d, start_port=%d, end_port=%d) on policy %s: no matching entry found in the policy's l4_params", params.Proto, params.StartPort, params.EndPort, policy_id)
}

// DescribePolicyL4Params fetches the parent policy and returns its raw id
// plus the list of l4_params (ports) attached to it. This is used both by
// CreatePort (to resolve the real l4_param id from the API's
// parent-policy-shaped response) and by the port resource's Read to detect
// drift.
func (c Client) DescribePolicyL4Params(policy_id string) (policyWithL4Params, error) {
	var resp policyWithL4Params
	policyURL := c.Config.APIURL + SecureWorkloadAPIV1BasePath + "/policies" + fmt.Sprintf("/%s", policy_id)
	request, err := signer.CreateJSONRequest(http.MethodGet, policyURL, nil)
	if err != nil {
		return resp, err
	}
	err = c.Do(request, &resp)
	return resp, err
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
	request, err := signer.CreateJSONRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return c.Do(request, nil)
}
