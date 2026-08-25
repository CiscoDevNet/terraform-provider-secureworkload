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
	// StartPort, EndPort, and Proto are pointers so that an unset value
	// can be omitted from the request body entirely (json:",omitempty" on
	// a plain int cannot distinguish "unset" from literal 0, and IANA
	// protocol/port numbers never legitimately use 0 anyway). All three
	// nil produces a bare "{}" body, which is exactly what CSW requires
	// to create an ANY service (all protocols, all ports): proto:null and
	// no "port" key at all. A port range MAY NOT be combined with a nil
	// proto -- CSW rejects that combination with "ports not enforceable
	// with wildcard proto" -- so callers must set all three or none.
	StartPort   *int   `json:"start_port,omitempty"`
	EndPort     *int   `json:"end_port,omitempty"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Proto       *int   `json:"proto,omitempty"`
}

// optionalIntString formats a *int for error messages, rendering nil as
// "unset" instead of a confusing zero value.
func optionalIntString(v *int) string {
	if v == nil {
		return "unset"
	}
	return fmt.Sprintf("%d", *v)
}

// isAnyRequest reports whether params describes CSW's ANY service (all
// protocols, all ports): no start_port, no end_port, and no proto.
func isAnyRequest(params CreatePortRequest) bool {
	return params.StartPort == nil && params.EndPort == nil && params.Proto == nil
}

// isAnyL4Param reports whether an already-normalize()'d l4_param is CSW's
// ANY service. The API represents ANY with proto:null and no "port" key at
// all; after JSON decoding that means no PortRange, and Proto/StartPort/
// EndPort all at their zero value, since IANA protocol/port numbers never
// legitimately use 0.
func isAnyL4Param(p Port) bool {
	return len(p.PortRange) == 0 && p.StartPort == 0 && p.EndPort == 0 && p.Proto == 0
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
// matches the requested proto/start_port/end_port. An ANY request (all
// three unset) only matches an ANY l4_param (proto:null, no "port" key);
// a concrete port/proto request only matches a concrete l4_param with the
// same values. Returns the matching Port and true if found.
func findMatchingL4Param(l4Params []Port, params CreatePortRequest) (Port, bool) {
	var match Port
	found := false
	wantAny := isAnyRequest(params)
	for i := range l4Params {
		p := l4Params[i]
		p.normalize()
		if wantAny {
			if isAnyL4Param(p) {
				match = p
				found = true
			}
			continue
		}
		if isAnyL4Param(p) {
			continue
		}
		if params.StartPort != nil && p.StartPort != *params.StartPort {
			continue
		}
		if params.EndPort != nil && p.EndPort != *params.EndPort {
			continue
		}
		if params.Proto != nil && p.Proto != *params.Proto {
			continue
		}
		match = p
		found = true
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

	return port, fmt.Errorf("could not determine the l4_param id for the newly created port (proto=%s, start_port=%s, end_port=%s) on policy %s: no matching entry found in the policy's l4_params", optionalIntString(params.Proto), optionalIntString(params.StartPort), optionalIntString(params.EndPort), policy_id)
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
