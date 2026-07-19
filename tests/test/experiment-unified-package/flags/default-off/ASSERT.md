## Expected

- Parse succeeds.
- `opts.ExperimentUnifiedPackagePerDoctestTree == false`.
- `opts.ExperimentRefInsteadOfInline == false`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse error: %s", resp.ParseErr)
	}
	if resp.Opts.ExperimentUnifiedPackagePerDoctestTree {
		t.Fatal("ExperimentUnifiedPackagePerDoctestTree=true by default; want false")
	}
	if resp.Opts.ExperimentRefInsteadOfInline {
		t.Fatal("ExperimentRefInsteadOfInline=true by default; want false (no experiment flags)")
	}
}
```
