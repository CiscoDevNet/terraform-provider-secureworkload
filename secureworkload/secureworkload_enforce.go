package secureworkload

import (
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
