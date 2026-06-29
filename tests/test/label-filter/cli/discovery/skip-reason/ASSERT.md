## Expected

- Skip block lists five leaves; each entry has `reason: label filter`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, resp.Stdout)
	}
	block := skipBlock(resp.Stdout)
	if block == "" {
		t.Fatal("missing skip block")
	}
	if strings.Count(block, "reason: label filter") != 5 {
		t.Fatalf("expected 5 reason lines, block:\n%s", block)
	}
	assertNoResultSummary(t, resp.Stdout)
}
```