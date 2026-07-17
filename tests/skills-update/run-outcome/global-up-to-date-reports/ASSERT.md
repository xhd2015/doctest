---
label: heavy
---

## Expected

- Exit code 0.
- stdout contains `Skill is up to date` for the global `doctest-tdd` install path.
- stdout contains `skill not installed:` for every registry skill except `tdd`.
- stdout does not contain `No installed skills found`.

## Expected Output

```
<contains>
Skill is up to date
doctest-tdd
skill not installed: code-spec
skill not installed: tdd-lite
</contains>
```

## Errors

- None from `Run`.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("expected up-to-date line:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "doctest-tdd") {
		t.Fatalf("stdout should reference doctest-tdd:\n%s", resp.Stdout)
	}
	for _, name := range registryCLINames() {
		if name == "tdd" {
			continue
		}
		assertNotInstalledLines(t, resp.Stdout, name)
	}
	assertNoScopeHint(t, resp.Stdout)
	assert.Output(t, resp.Stdout, `` +
`<contains>
Skill is up to date
doctest-tdd
skill not installed: code-spec
skill not installed: tdd-lite
</contains>`)
}
```
