---
label: heavy
explanation: CLI filter contract via doctest binary
---

## Expected

- PASS(1/1); no skip block.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, resp.Stdout)
	}
	assertNoSkipBlock(t, resp.Stdout)
	assertResultSummary(t, resp.Stdout, 1, 1)
}
```