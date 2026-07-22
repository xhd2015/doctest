---
label: heavy
---

## Expected

- Exit code 0.
- stdout contains `skill not installed: tdd` (global install invisible without `--global`).
- stdout contains not-installed lines for every other registry skill.
- stdout does not contain `Skill is up to date`, `Update skill at`, or `No installed skills found`.

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

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	assertNotInstalledLines(t, resp.Stdout, registryCLINames()...)
	assertNoScopeHint(t, resp.Stdout)
	if strings.Contains(resp.Stdout, "Skill is up to date") || strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("expected no per-skill update lines without --global:\n%s", resp.Stdout)
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