---
explanation: CLI filter contract via doctest binary
---

## Expected

- PASS(2/2); compact skip buckets for non-matching leaves.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stdout:\n%s stderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertLabelFilterSkipCompact(t, resp.Stdout, 3, map[string]int{
		"(unlabeled)":   1,
		"flaky":         1,
		"ui-automation": 1,
	})
	assertResultSummary(t, resp.Stdout, 2, 2)
}
```
