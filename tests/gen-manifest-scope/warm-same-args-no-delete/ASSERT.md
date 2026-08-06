---
explanation: warm same-args; gen-plan deleted=0; tree stays managed
---

## Expected

- Both runs exit 0.
- After run1, manifest lists `tree/` packages.
- After run2 (identical args), manifest still lists `tree/`.
- Run2 gen-plan summary has `deleted=0` (no managed orphans on warm).
- Prefer run2 also shows some `unchanged` activity when gen-plan is on
  (soft: only require deleted=0 if summary parse is partial).

## Errors

- Fail if warm run reports deleted>0.
- Fail if tree ledger entries disappear.

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
		t.Fatalf("warm run exit=%d\nstderr:\n%s", resp.SubsetExitCode, resp.SubsetStderr)
	}
	if resp.ManifestAfterFull == "" {
		t.Fatalf("missing manifest after cold under %s", req.GenDir)
	}
	if !manifestHasTreePrefix(resp.ManifestAfterFull, "tree") {
		t.Fatalf("cold manifest missing tree/ entries:\n%s", resp.ManifestAfterFull)
	}
	if !manifestHasTreePrefix(resp.ManifestAfterSubset, "tree") {
		t.Fatalf("warm manifest lost tree/ entries:\n%s", resp.ManifestAfterSubset)
	}

	del := parseDeletedFromGenPlan(resp.SubsetStderr)
	if del < 0 {
		// gen-plan may be absent if debug key missing — still require no delete signal.
		if strings.Contains(resp.SubsetStderr, "# deleted") {
			t.Fatalf("warm stderr mentions # deleted without parseable summary:\n%s", resp.SubsetStderr)
		}
		return
	}
	if del != 0 {
		t.Fatalf("warm same-args expected deleted=0, got deleted=%d\nstderr:\n%s", del, resp.SubsetStderr)
	}
}
```
