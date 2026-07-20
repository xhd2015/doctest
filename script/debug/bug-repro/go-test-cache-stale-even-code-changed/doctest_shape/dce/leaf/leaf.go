package leaf

import (
	"testing"

	"dt.local/droot"
	"dt.local/mid"
	"dt.local/registry"
)

func init() {
	registry.Register(registry.Entry{Path: "leaf", Fn: RunTestLeaf})
}

// WorkDir is written by mid.Setup and never read (empty assert).
// The compiler may drop that write from the linked test binary (DCE),
// so binary content ID can stay unchanged after editing mid.go.
func RunTestLeaf(t *testing.T) {
	req := &droot.Request{}
	_ = droot.RootSetup(t, req)
	_ = mid.Setup(t, req)
	_, _ = droot.Run(t, req)
}
