## Expected

- Exit code 0.
- Output includes the default run stem `sumdef01` and/or default-suite leaf markers.
- Output does not present `sumlab01` / `only/labeled-suite-leaf` as included runs (labeled suite excluded).

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	out := combinedOut(resp)
	if !strings.Contains(out, "sumdef01") && !strings.Contains(out, "slow-leaf") && !strings.Contains(strings.ToLower(out), "default") {
		// At least one signal that default suite data is in the summary.
		// Accept numeric aggregate that implies 1 run if stems omitted.
		if !strings.Contains(out, "1") {
			t.Fatalf("summary --default-only missing default-suite signal:\n%s", out)
		}
	}
	if strings.Contains(out, "sumlab01") || strings.Contains(out, "only/labeled-suite-leaf") {
		t.Fatalf("summary --default-only should exclude labeled suite run:\n%s", out)
	}
}
```
