## Expected

- Parse succeeds.
- `opts.MetricsOn == true`.
- Remain args include `./tests`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse error: %s", resp.ParseErr)
	}
	if !resp.Opts.MetricsOn {
		t.Fatal("MetricsOn=false after --metrics-on; want true")
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
