## Expected

- Suite run succeeds (`RunErr` empty) even when `ExperimentRefInsteadOfInline` is true.
- No harness error.
- (P0: does not assert ref-based generated source.)

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("RunTest failed with ExperimentRefInsteadOfInline=true (P0 should still use classic gen): %s\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, resp.Stdout, resp.Stderr)
	}
}
```
