## Expected

- Exit code 0.
- stdout contains polished `tdd  up to date` (global install under
  `~/.agents/skills/doctest-tdd`, display name is CLI `tdd`).
- stdout contains `<name>  not installed` for every other registry skill.
- summary includes `1 up to date` and `14 not installed`.
- stdout does not contain `No installed skills found`.

## Expected Output

```
<contains>
tdd  up to date
code-spec  not installed
tdd-lite  not installed
0 updated · 1 up to date · 14 not installed
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

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	assertUpToDateLine(t, resp.Stdout, "tdd")
	for _, name := range registryCLINames() {
		if name == "tdd" {
			continue
		}
		assertNotInstalledLines(t, resp.Stdout, name)
	}
	assertNoScopeHint(t, resp.Stdout)
	plain := stripANSI(resp.Stdout)
	if !strings.Contains(plain, "0 updated · 1 up to date · 14 not installed") {
		t.Fatalf("stdout missing batch summary:\n%s", resp.Stdout)
	}
	assert.Output(t, plain, ``+
		`<contains>
tdd  up to date
code-spec  not installed
tdd-lite  not installed
0 updated · 1 up to date · 14 not installed
</contains>`)
}
```
