---
label: heavy
explanation: source change → gen # modified
---

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	requireExit0(t, "run1", resp.FullExitCode, resp.FullStderr)
	requireExit0(t, "run2", resp.SubsetExitCode, resp.SubsetStderr)

	if !manifestHasExactRel(resp.ManifestAfterSubset, "tree/leaf/leaf.go") {
		t.Fatalf("leaf.go missing from manifest after source change:\n%s", resp.ManifestAfterSubset)
	}
	// Prefer explicit # modified on leaf.go; fallback to summary modified>=1.
	if !genPlanHasTag(resp.SubsetStderr, "leaf.go", "modified") {
		mod := parseModifiedFromGenPlan(resp.SubsetStderr)
		if mod < 1 {
			// Some paths may rewrite as new if package path changes — accept new|modified.
			if !genPlanHasTag(resp.SubsetStderr, "leaf.go", "new") && parseNewFromGenPlan(resp.SubsetStderr) < 1 {
				t.Fatalf("expected # modified (or new) after source change:\n%s", resp.SubsetStderr)
			}
		}
	}
}
```
