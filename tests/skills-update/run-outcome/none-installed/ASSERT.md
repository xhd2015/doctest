---
label: heavy
---

## Expected

- Exit code 0.
- stdout contains `skill not installed: <name>` for every registry CLI skill.
- stdout does not contain `No installed skills found`.

## Expected Output

```
<contains>
skill not installed: analyse-perf
skill not installed: code-spec
skill not installed: design-principle
skill not installed: designer
skill not installed: doc-spec
skill not installed: implementer
skill not installed: lint
skill not installed: migrate
skill not installed: output-assert
skill not installed: reproduce
skill not installed: review
skill not installed: review-perf
skill not installed: tdd
skill not installed: tdd-cli-agent
skill not installed: tdd-lite
</contains>
```

## Errors

- `Run` returns no error.

```go
import (
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
}
```