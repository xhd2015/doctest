---
explanation: removed source leaf → gen orphan # deleted
---

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	requireExit0(t, "run1", resp.FullExitCode, resp.FullStderr)
	requireExit0(t, "run2", resp.SubsetExitCode, resp.SubsetStderr)

	goneRel := "tree/gone/leaf.go"
	if !manifestHasExactRel(resp.ManifestAfterFull, goneRel) {
		t.Fatalf("cold manifest missing %s:\n%s", goneRel, resp.ManifestAfterFull)
	}
	// After prune: not on disk, not in ledger.
	goneAbs := filepath.Join(req.GenDir, filepath.FromSlash(goneRel))
	if _, err := os.Stat(goneAbs); err == nil {
		t.Fatalf("orphan gen file still on disk: %s", goneAbs)
	}
	if manifestHasExactRel(resp.ManifestAfterSubset, goneRel) {
		t.Fatalf("orphan still in manifest after prune:\n%s", resp.ManifestAfterSubset)
	}
	// keep leaf still managed.
	if !manifestHasExactRel(resp.ManifestAfterSubset, "tree/keep/leaf.go") {
		t.Fatalf("keep leaf must remain in manifest:\n%s", resp.ManifestAfterSubset)
	}
	if !genPlanHasTag(resp.SubsetStderr, "leaf.go", "deleted") {
		del := parseDeletedFromGenPlan(resp.SubsetStderr)
		if del < 1 {
			t.Fatalf("expected # deleted / deleted>=1 for removed leaf; stderr:\n%s", resp.SubsetStderr)
		}
	}
}
```
