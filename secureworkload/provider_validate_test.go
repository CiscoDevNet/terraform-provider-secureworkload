package secureworkload

import "testing"

// TestProviderInternalValidate guards against illegal schema attribute
// combinations (e.g. Optional+ForceNew+non-nil Update on a resource with
// no other updatable fields, or Computed+Default on the same attribute).
// It is intentionally kept as a permanent regression test: the SDK only
// catches these classes of bugs via InternalValidate, not via `go build`
// or `go vet`.
func TestProviderInternalValidate(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("provider internal validate: %s", err)
	}
}
