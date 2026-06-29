## Expected

- Exit 0; slow leaf skipped; no PASS line.

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, resp.Stdout)
	}
	mod := filepath.Dir(req.Args[1])
	want := wantLabelFilterSkipBlockMulti(1,
		wantLabelFilterSkipEntry(mod, "slow", "slow", "slow profile", false),
	)
	assertSkipBlockExact(t, resp.Stdout, want)
	assertNoResultSummary(t, resp.Stdout)
}
```