## Expected

- Exit code 0.
- stdout contains polished `<name>  not installed` for every registry CLI skill
  (skills v0.0.26+ batch format).
- stdout ends with a summary including `not installed`.
- stdout does not contain `No installed skills found`.

## Expected Output

```
<contains>
analyse-perf  not installed
code-spec  not installed
design-principle  not installed
designer  not installed
doc-spec  not installed
implementer  not installed
lint  not installed
migrate  not installed
output-assert  not installed
reproduce  not installed
review  not installed
review-perf  not installed
tdd  not installed
tdd-cli-agent  not installed
tdd-lite  not installed
0 updated · 0 up to date · 15 not installed
</contains>
```

## Errors

- `Run` returns no error.

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
		t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	assertNotInstalledLines(t, resp.Stdout, registryCLINames()...)
	assertNoScopeHint(t, resp.Stdout)
	plain := stripANSI(resp.Stdout)
	if !strings.Contains(plain, "0 updated · 0 up to date · 15 not installed") {
		t.Fatalf("stdout missing batch summary:\n%s", resp.Stdout)
	}
}
```