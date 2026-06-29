## Expected

- Exit 0; all five leaves skipped; no PASS/FAIL summary.

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
	want := wantLabelFilterSkipBlockMulti(5,
		wantLabelFilterSkipEntry(mod, "both", "slow, ui-automation", "slow ui combo", false),
		wantLabelFilterSkipEntry(mod, "fast", "", "", true),
		wantLabelFilterSkipEntry(mod, "heavy", "heavy", "heavy profile", false),
		wantLabelFilterSkipEntry(mod, "slow", "slow", "slow profile", false),
		wantLabelFilterSkipEntry(mod, "ui", "ui-automation", "browser ui", false),
	)
	assertSkipBlockExact(t, resp.Stdout, want)
	assertNoResultSummary(t, resp.Stdout)
}
```