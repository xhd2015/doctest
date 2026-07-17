---
label: heavy
---

## Expected

- Same outcome as OR expression: PASS(3/3), skips fast and ui.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, resp.Stdout)
	}
	mod := req.Args[1]
	want := wantLabelFilterSkipBlockMulti(2,
		wantLabelFilterSkipEntry(mod, "fast", "", "", true),
		wantLabelFilterSkipEntry(mod, "ui", "ui-automation", "browser ui", false),
	)
	assertSkipBlockExact(t, resp.Stdout, want)
	assertResultSummary(t, resp.Stdout, 3, 3)
}
```