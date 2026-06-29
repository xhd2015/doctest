## Expected

- Exit code 0.
- stdout contains `Skill is up to date` referencing `doctest-tdd`.
- stdout contains `skill not installed:` for every registry skill except `tdd`.
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

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Skill is up to date") && !strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("expected update output for installed skill:\n%s", resp.Stdout)
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
	skillPath := filepath.Join(req.WorkDir, ".agents", "skills", "doctest-tdd", "SKILL.md")
	if _, err := os.Stat(skillPath); err != nil {
		t.Fatalf("expected installed skill at %s: %v", skillPath, err)
	}
}
```