---
label: heavy
---

## Expected

- Compact skip header indicates label filter; five leaves bucketed by label set.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, resp.Stdout)
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
