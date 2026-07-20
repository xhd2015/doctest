---
label: heavy
---

## Expected

- Exit 0; single compact skip for slow leaf (label filter miss).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stdout:\n%s stderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertLabelFilterSkipCompact(t, resp.Stdout, 1, map[string]int{
		"slow": 1,
	})
	assertNoResultSummary(t, resp.Stdout)
}
```
