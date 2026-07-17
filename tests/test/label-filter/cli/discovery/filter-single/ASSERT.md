---
label: heavy
---

## Expected

- PASS(2/2); skips fast, ui, heavy with label-filter reason.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stdout:\n%s stderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	mod := req.Args[1]
	want := wantLabelFilterSkipBlockMulti(3,
		wantLabelFilterSkipEntry(mod, "fast", "", "", true),
		wantLabelFilterSkipEntry(mod, "heavy", "heavy", "heavy profile", false),
		wantLabelFilterSkipEntry(mod, "ui", "ui-automation", "browser ui", false),
	)
	assertSkipBlockExact(t, resp.Stdout, want)
	assertResultSummary(t, resp.Stdout, 2, 2)
}
```