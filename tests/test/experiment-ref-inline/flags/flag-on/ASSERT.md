## Expected

- Parse succeeds.
- `opts.ExperimentRefInsteadOfInline == true`.
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
	if !resp.Opts.ExperimentRefInsteadOfInline {
		t.Fatal("ExperimentRefInsteadOfInline=false after --experiment-ref-instead-of-inline; want true")
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
