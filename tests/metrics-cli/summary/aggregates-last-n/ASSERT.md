## Expected

- Exit code 0.
- Output mentions both newest stems (`summid01`, `sumnew01`) **or** reports an aggregate run count of 2 (or "2 runs").
- Output should not require the oldest stem `sumold01` (may omit it).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := combinedOut(resp)
	hasMid := strings.Contains(out, "summid01")
	hasNew := strings.Contains(out, "sumnew01")
	hasCount2 := strings.Contains(out, "2") && (strings.Contains(strings.ToLower(out), "run"))
	if !(hasMid && hasNew) && !hasCount2 {
		t.Fatalf("summary --last 2 should cover two newest runs:\n%s", out)
	}
}
```
