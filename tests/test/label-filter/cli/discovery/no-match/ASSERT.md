---
label: heavy
---

## Expected

- Exit 0; all five leaves in compact label-filter skip; no PASS line.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stdout:\n%s stderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertLabelFilterSkipCompact(t, resp.Stdout, 5, map[string]int{
		"(unlabeled)":        1,
		"heavy":              1,
		"slow":               1,
		"slow,ui-automation": 1,
		"ui-automation":      1,
	})
	assertNoResultSummary(t, resp.Stdout)
}
```
