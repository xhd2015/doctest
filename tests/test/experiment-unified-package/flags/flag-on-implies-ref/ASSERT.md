## Expected

- Parse succeeds.
- `opts.ExperimentUnifiedPackagePerDoctestTree == true`.
- `opts.ExperimentRefInsteadOfInline == true` (auto-enabled by unified flag).
- Remain args include `./tests`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse error: %s", resp.ParseErr)
	}
	if !resp.Opts.ExperimentUnifiedPackagePerDoctestTree {
		t.Fatal("ExperimentUnifiedPackagePerDoctestTree=false after --experiment-unified-package-per-doctest-tree; want true")
	}
	if !resp.Opts.ExperimentRefInsteadOfInline {
		t.Fatal("ExperimentRefInsteadOfInline=false after unified flag; want true (unified auto-enables ref)")
	}
	found := false
	for _, a := range resp.RemainArgs {
		if a == "./tests" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("remain args missing ./tests: %v", resp.RemainArgs)
	}
}
```
