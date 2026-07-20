---
label: heavy
---

## Expected

- PASS(1/1); four compact skip buckets.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stdout:\n%s stderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertLabelFilterSkipCompact(t, resp.Stdout, 4, map[string]int{
		"(unlabeled)":   1,
		"heavy":         1,
		"slow":          1,
		"ui-automation": 1,
	})
	assertResultSummary(t, resp.Stdout, 1, 1)
}
```
