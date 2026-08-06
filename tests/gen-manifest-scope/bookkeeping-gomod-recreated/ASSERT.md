---
explanation: missing go.mod recreated on next generate
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

	if !resp.GoModExistsAfter2 {
		t.Fatal("go.mod not recreated after delete")
	}
	if !manifestHasExactRel(resp.ManifestAfterSubset, "go.mod") {
		t.Fatalf("go.mod missing from manifest after recreate:\n%s", resp.ManifestAfterSubset)
	}
}
```
