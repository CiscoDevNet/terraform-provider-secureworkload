package secureworkload

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// ApplicationVersion represents a single enforceable policy version of a workspace,
// as returned by the `/applications/{id}/versions` endpoint.
type ApplicationVersion struct {
	Version                string `json:"version"`
	CreatedAt              int64  `json:"created_at"`
	Name                   string `json:"name,omitempty"`
	Description            string `json:"description,omitempty"`
	LastEnforcementEventAt int64  `json:"last_enforcement_event_at,omitempty"`
}

// isNotModified reports whether err is (or wraps) an *APIError with a 304
// Not Modified status code. For an idempotent enable_enforce call, a 304
// means enforcement is already active at the requested version -- which is
// exactly the success condition, not a failure.
func isNotModified(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotModified
	}
	return false
}

// isAlreadyEnforced reports whether err is the API's way of saying the
// requested version is already the enforced one, e.g.
//
//	400  {"error": "version p76 is already enforced", "status": "bad_request"}
//
// enable_enforce is idempotent in intent: if the requested version is
// already live, the desired end state is satisfied and there is nothing to
// do. Without this, adopting a workspace that is already enforcing is
// impossible, because Create (and Update, when track_latest_version
// resolves to the version already enforced) fails outright.
//
// The match is deliberately narrow: a 400 alone is NOT enough, since a bad
// version string or a malformed body also returns 400 and those must keep
// failing. Only a 400 whose body actually says "already enforced" is
// treated as success.
//
// Caveat: this matches on the API's error text, so a future rewording by
// Cisco would silently stop it matching. That fails safe (an error, not
// silent corruption) but it is not robust, and a dedicated
// "get enforcement status" endpoint would be a better long-term answer.
func isAlreadyEnforced(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		return false
	}
	body, mErr := json.Marshal(apiErr.Body)
	if mErr != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(body)), "already enforced")
}

func (c Client) CreateEnforce(params CreateEnforceRequest, workspace_id string) (Enforce, error) {
	var enforce Enforce
	url := c.Config.APIURL + EnforceAPIV1BasePath + workspace_id + "/enable_enforce"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, params)
	if err != nil {
		return enforce, err
	}
	err = c.Do(request, &enforce)
	if err != nil && (isNotModified(err) || isAlreadyEnforced(err)) {
		// Already enforced at the requested version -- treat as success.
		// 304 covers the "no version supplied" path; 400 "already
		// enforced" covers the explicit-version path that
		// track_latest_version takes.
		return enforce, nil
	}
	return enforce, err
}

// ListApplicationVersions returns the list of enforceable policy ("p<N>") versions
// for a workspace, newest first (per the API), though callers should not rely on
// ordering and should compute the max explicitly.
func (c Client) ListApplicationVersions(workspace_id string) ([]ApplicationVersion, error) {
	var versions []ApplicationVersion
	url := c.Config.APIURL + EnforceAPIV1BasePath + workspace_id + "/versions"
	request, err := signer.CreateJSONRequest(http.MethodGet, url, nil)
	if err != nil {
		return versions, err
	}
	err = c.Do(request, &versions)
	return versions, err
}

// parsePolicyVersion parses a policy version string such as "p73" (or a bare "73")
// into its integer component. It rejects "v"-prefixed strings (ADM versions), since
// those are a different counter from enforceable policy versions.
func parsePolicyVersion(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if strings.HasPrefix(s, "v") || strings.HasPrefix(s, "V") {
		return 0, false
	}
	numStr := s
	if strings.HasPrefix(s, "p") || strings.HasPrefix(s, "P") {
		numStr = s[1:]
	}
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, false
	}
	return n, true
}

// formatPolicyVersion formats an integer policy version as "p<N>".
func formatPolicyVersion(n int) string {
	return fmt.Sprintf("p%d", n)
}

// LatestPolicyVersion returns the newest enforceable policy version ("p<N>") for a
// workspace, computed by explicitly finding the max parsed version rather than
// relying on the list's ordering.
func (c Client) LatestPolicyVersion(workspace_id string) (string, error) {
	versions, err := c.ListApplicationVersions(workspace_id)
	if err != nil {
		return "", err
	}
	best := -1
	found := false
	for _, v := range versions {
		if n, ok := parsePolicyVersion(v.Version); ok {
			if !found || n > best {
				best = n
				found = true
			}
		}
	}
	if !found {
		return "", fmt.Errorf("no enforceable policy versions found for workspace %s", workspace_id)
	}
	return formatPolicyVersion(best), nil
}

func (c Client) DeleteEnforce(workspace_id string) error {
	url := c.Config.APIURL + EnforceAPIV1BasePath + workspace_id + "/disable_enforce"
	request, err := signer.CreateJSONRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	if err := c.Do(request, nil); err != nil {
		return err
	}
	// Poll for enforcement to actually turn off, bounded at ~60s total.
	for i := 0; i < 12; i++ {
		time.Sleep(5 * time.Second)
		application, err := c.DescribeApplication(DescribeApplicationRequest{ApplicationId: workspace_id})
		if err != nil {
			if IsNotFound(err) {
				return nil
			}
			return err
		}
		if !application.EnforcementEnabled {
			return nil
		}
	}
	return fmt.Errorf("timed out waiting for enforcement to be disabled on workspace %s", workspace_id)
}
