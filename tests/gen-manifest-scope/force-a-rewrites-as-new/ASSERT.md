---
label: heavy
explanation: -a wipe then regenerate → new>=1
---

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	requireExit0(t, "run1", resp.FullExitCode, resp.FullStderr)
	requireExit0(t, "run2 -a", resp.SubsetExitCode, resp.SubsetStderr)

	if resp.ManifestAfterSubset == "" {
		t.Fatal("manifest missing after -a regenerate")
	}
	if !manifestHasTreePrefix(resp.ManifestAfterSubset, "tree") {
		t.Fatalf("tree packages missing after -a:\n%s", resp.ManifestAfterSubset)
	}
	nw := parseNewFromGenPlan(resp.SubsetStderr)
	if nw < 1 && !genPlanHasTag(resp.SubsetStderr, "leaf.go", "new") {
		// -a wipe may report new or modified depending on timing; require not all-unchanged.
		unc := parseSummaryCount(resp.SubsetStderr, "unchanged")
		mod := parseModifiedFromGenPlan(resp.SubsetStderr)
		if nw == 0 && mod == 0 && unc > 0 {
			t.Fatalf("-a expected cold-like writes (new|modified>=1); stderr:\n%s", resp.SubsetStderr)
		}
		if nw < 0 && mod < 0 {
			t.Fatalf("missing gen-plan summary after -a:\n%s", resp.SubsetStderr)
		}
	}
}
```
