---
explanation: unmanaged plant under gen package survives; not in manifest
---

## Expected

- Both runs exit 0.
- After run2, planted file still exists on disk.
- Manifest does **not** list `tree/__droot/unused.go` (never managed).
- gen-plan warm/second summary: `deleted=0` (plant is not a ledger orphan).
- Second stderr does not show `unused.go` with `# deleted`.

## Errors

- Fail if plant was removed (disk-walk orphan policy).
- Fail if plant was adopted into the manifest without a real write path.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run error: %v\nrun1 stderr:\n%s\nrun2 stderr:\n%s", err, resp.FullStderr, resp.SubsetStderr)
	}
	if resp.FullExitCode != 0 {
		t.Fatalf("cold run exit=%d\nstderr:\n%s", resp.FullExitCode, resp.FullStderr)
	}
	if resp.SubsetExitCode != 0 {
		t.Fatalf("second run exit=%d\nstderr:\n%s", resp.SubsetExitCode, resp.SubsetStderr)
	}
	if resp.PlantAbs == "" {
		t.Fatal("expected PlantAbs to be set when PlantRel is configured")
	}
	if !resp.PlantExistsAfter {
		t.Fatalf("unmanaged plant was deleted (must survive; not in manifest): %s", resp.PlantAbs)
	}
	if resp.PlantInManifestAfter {
		t.Fatalf("unmanaged plant must not appear in %s:\n%s", genManifestName, resp.ManifestAfterSubset)
	}
	// gen-plan: not listed as deleted
	if strings.Contains(resp.SubsetStderr, "unused.go") && strings.Contains(resp.SubsetStderr, "# deleted") {
		// only fail if unused.go line carries the deleted tag nearby
		for _, line := range strings.Split(resp.SubsetStderr, "\n") {
			if strings.Contains(line, "unused.go") && strings.Contains(line, "# deleted") {
				t.Fatalf("plant must not appear as # deleted:\n%s\nfull stderr:\n%s", line, resp.SubsetStderr)
			}
		}
	}
	del := parseDeletedFromGenPlan(resp.SubsetStderr)
	if del > 0 {
		t.Fatalf("second run deleted=%d; unmanaged plant must not create ledger deletes\nstderr:\n%s", del, resp.SubsetStderr)
	}
}
```
