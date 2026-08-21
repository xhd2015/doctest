## Expected

- Exit code 0.
- stdout contains `<name>  not installed` for every registry skill including
  `tdd` (global install invisible without `--global`).
- stdout does **not** contain `tdd  up to date` or legacy `Skill is up to date` /
  `Update skill at`.
- summary is `0 updated · 0 up to date · 16 not installed`.
- stdout does not contain `No installed skills found`.

## Side Effects

- Global install under `$HOME/.agents/skills/doctest-tdd` remains; no local
  `.agents/skills` tree is created.

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
	assertNotInstalledLines(t, resp.Stdout, registryCLINames()...)
	assertNoScopeHint(t, resp.Stdout)
	plain := stripANSI(resp.Stdout)
	if strings.Contains(plain, "tdd  up to date") ||
		strings.Contains(plain, "Skill is up to date") ||
		strings.Contains(plain, "Update skill at") {
		t.Fatalf("expected no up-to-date lines without --global:\n%s", resp.Stdout)
	}
	if !strings.Contains(plain, "0 updated · 0 up to date · 16 not installed") {
		t.Fatalf("stdout missing batch summary:\n%s", resp.Stdout)
	}
	if req.Home == "" {
		t.Fatal("req.Home unset (isolated HOME for global install)")
	}
	globalSkill := filepath.Join(req.Home, ".agents", "skills", "doctest-tdd", "SKILL.md")
	if _, err := os.Stat(globalSkill); err != nil {
		t.Fatalf("expected global install at %s: %v", globalSkill, err)
	}
	localSkill := filepath.Join(req.WorkDir, ".agents", "skills", "doctest-tdd", "SKILL.md")
	if _, err := os.Stat(localSkill); err == nil {
		t.Fatalf("expected no local install at %s", localSkill)
	}
}
```