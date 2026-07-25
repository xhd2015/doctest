---
label: heavy
explanation: managed missing gen file recreated as # new
---

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run: %v\n%s\n%s", err, resp.FullStderr, resp.SubsetStderr)
	}
	requireExit0(t, "run1", resp.FullExitCode, resp.FullStderr)
	requireExit0(t, "run2", resp.SubsetExitCode, resp.SubsetStderr)

	if !manifestHasExactRel(resp.ManifestAfterFull, "tree/leaf/leaf.go") {
		t.Fatalf("cold manifest missing tree/leaf/leaf.go:\n%s", resp.ManifestAfterFull)
	}
	if !resp.DeletedGenExistsAfter2 {
		t.Fatalf("managed file not recreated: %s", resp.DeletedGenAbs)
	}
	if !manifestHasExactRel(resp.ManifestAfterSubset, "tree/leaf/leaf.go") {
		t.Fatalf("manifest lost leaf.go after rewrite:\n%s", resp.ManifestAfterSubset)
	}
	if !genPlanHasTag(resp.SubsetStderr, "leaf.go", "new") {
		nw := parseNewFromGenPlan(resp.SubsetStderr)
		if nw < 1 {
			t.Fatalf("expected # new for recreated leaf.go (or new>=1); stderr:\n%s", resp.SubsetStderr)
		}
	}
}
```
