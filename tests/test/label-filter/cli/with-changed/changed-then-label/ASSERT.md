## Expected

- Only changed `slow` leaf is eligible; label filter matches; PASS(1/1).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stdout:\n%s stderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	assertNoSkipBlock(t, resp.Stdout)
	assertResultSummary(t, resp.Stdout, 1, 1)
}
```