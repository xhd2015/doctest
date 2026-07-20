package leaf

import (
	"testing"

	"dt.local/droot"
	"dt.local/mid"
	"dt.local/registry"
)

func init() {
	registry.Register(registry.Entry{Path: "mid-a/leaf-x", Fn: RunTestLeaf})
}

func RunTestLeaf(t *testing.T) {
	req := &droot.Request{}
	_ = droot.RootSetup(t, req)
	_ = mid.Setup(t, req)
	_, _ = droot.Run(t, req)
	if req.WorkDir != "MID_V1" {
		t.Fatalf("WorkDir=%q want MID_V1", req.WorkDir)
	}
}
