---
label: heavy
---

## Expected

- Exit code 0.
- stdout contains polished `tdd  up to date` (CLI registry name; install dir remains
  `doctest-tdd`).
- stdout contains `<name>  not installed` for every other registry skill.
- summary includes `1 up to date` and `14 not installed`.
- stdout does not contain `No installed skills found`.

## Side Effects

- `SKILL.md` exists under `.agents/skills/doctest-tdd` in `WorkDir`.

## Errors

- None from `Run`.

```go
import (
	"os"
	"path/filepath"
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
	skillPath := filepath.Join(req.WorkDir, ".agents", "skills", "doctest-tdd", "SKILL.md")
	if _, err := os.Stat(skillPath); err != nil {
		t.Fatalf("expected installed skill at %s: %v", skillPath, err)
	}
}
```