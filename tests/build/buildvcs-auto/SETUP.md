# Scenario

**Feature**: shared origin-module + clone helpers for -buildvcs=auto leaves

```
# origin module (multi-commit) is the only clone source
createOriginModule -> file:// origin -> clone full | shallow | shallow+break HEAD

# each leaf builds in its clone with -buildvcs=auto
go build -buildvcs=auto -o app.bin . -> exit 0 | error obtaining VCS status
```

## Preconditions

- `git` and `go` are on PATH.
- Leaves create isolated temp dirs (`t.TempDir`); no dependence on the host
  worktree being shallow or full.

## Steps

1. Provide helpers to create a multi-commit origin module.
2. Provide full clone, shallow clone, and broken-HEAD clone helpers.
3. Leaf Setups call a helper and set `req.CloneDir`.

## Context

- Commits use `core.hooksPath=/dev/null` so global hooks cannot block the suite.
- Origin uses `main` as initial branch when supported.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	_ = req
	return nil
}

func gitCmd(dir string, args ...string) *exec.Cmd {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=buildvcs-auto",
		"GIT_AUTHOR_EMAIL=buildvcs-auto@example.com",
		"GIT_COMMITTER_NAME=buildvcs-auto",
		"GIT_COMMITTER_EMAIL=buildvcs-auto@example.com",
	)
	return cmd
}

func mustGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	out, err := gitCmd(dir, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// createOriginModule writes a minimal main module with two commits and returns
// the absolute path of the bare-ready origin working tree.
func createOriginModule(t *testing.T) string {
	t.Helper()
	origin := filepath.Join(t.TempDir(), "origin")
	if err := os.MkdirAll(origin, 0755); err != nil {
		t.Fatalf("mkdir origin: %v", err)
	}
	if err := os.WriteFile(filepath.Join(origin, "go.mod"), []byte("module example.com/buildvcsapp\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	main1 := "package main\n\nimport \"fmt\"\n\n" + "func main() { fmt.Println(\"v1\") }\n"
	if err := os.WriteFile(filepath.Join(origin, "main.go"), []byte(main1), 0644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	// Prefer -b main; fall back if older git.
	if out, err := gitCmd(origin, "init", "-b", "main").CombinedOutput(); err != nil {
		if out2, err2 := gitCmd(origin, "init").CombinedOutput(); err2 != nil {
			t.Fatalf("git init: %v\n%s\nfallback: %v\n%s", err, out, err2, out2)
		}
	}
	mustGit(t, origin, "config", "core.hooksPath", "/dev/null")
	mustGit(t, origin, "config", "user.email", "buildvcs-auto@example.com")
	mustGit(t, origin, "config", "user.name", "buildvcs-auto")
	mustGit(t, origin, "add", "go.mod", "main.go")
	mustGit(t, origin, "commit", "-m", "init")

	main2 := "package main\n\nimport \"fmt\"\n\n" + "func main() { fmt.Println(\"v2\") }\n"
	if err := os.WriteFile(filepath.Join(origin, "main.go"), []byte(main2), 0644); err != nil {
		t.Fatalf("write main.go v2: %v", err)
	}
	mustGit(t, origin, "add", "main.go")
	mustGit(t, origin, "commit", "-m", "v2")

	// Third commit so shallow depth=1 is meaningfully shorter than full.
	main3 := "package main\n\nimport \"fmt\"\n\n" + "func main() { fmt.Println(\"v3\") }\n"
	if err := os.WriteFile(filepath.Join(origin, "main.go"), []byte(main3), 0644); err != nil {
		t.Fatalf("write main.go v3: %v", err)
	}
	mustGit(t, origin, "add", "main.go")
	mustGit(t, origin, "commit", "-m", "v3")
	return origin
}

func cloneOrigin(t *testing.T, origin string, destName string, shallow bool) string {
	t.Helper()
	dest := filepath.Join(t.TempDir(), destName)
	args := []string{"clone"}
	if shallow {
		args = append(args, "--depth", "1")
	}
	// file:// URL works for local paths on all platforms Go supports here.
	args = append(args, "file://"+origin, dest)
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git clone (shallow=%v): %v\n%s", shallow, err, out)
	}
	// Disable hooks in the clone as well.
	mustGit(t, dest, "config", "core.hooksPath", "/dev/null")
	return dest
}

func cloneFull(t *testing.T, origin string) string {
	t.Helper()
	return cloneOrigin(t, origin, "full", false)
}

func cloneShallow(t *testing.T, origin string) string {
	t.Helper()
	dest := cloneOrigin(t, origin, "shallow", true)
	out, err := gitCmd(dest, "rev-parse", "--is-shallow-repository").CombinedOutput()
	if err != nil {
		t.Fatalf("rev-parse shallow: %v\n%s", err, out)
	}
	if strings.TrimSpace(string(out)) != "true" {
		t.Fatalf("expected shallow clone, got %q", strings.TrimSpace(string(out)))
	}
	return dest
}

// cloneBrokenHEAD clones fully then corrupts .git/HEAD so git status exits 128.
// This models the real -buildvcs=auto failure path (not shallow depth).
func cloneBrokenHEAD(t *testing.T, origin string) string {
	t.Helper()
	dest := cloneFull(t, origin)
	headPath := filepath.Join(dest, ".git", "HEAD")
	if err := os.WriteFile(headPath, []byte("garbage-not-a-ref\n"), 0644); err != nil {
		t.Fatalf("corrupt HEAD: %v", err)
	}
	// Sanity: git status must fail (same class of failure Go surfaces).
	if out, err := gitCmd(dest, "status", "--porcelain").CombinedOutput(); err == nil {
		t.Fatalf("expected git status to fail after corrupt HEAD; out=%s", out)
	}
	return dest
}

// assertVCSStamped runs go version -m and requires vcs build settings.
func assertVCSStamped(t *testing.T, bin string) {
	t.Helper()
	out, err := exec.Command("go", "version", "-m", bin).CombinedOutput()
	if err != nil {
		t.Fatalf("go version -m: %v\n%s", err, out)
	}
	s := string(out)
	// go version -m lists settings as tab-separated lines including "vcs".
	if !strings.Contains(s, "vcs") {
		t.Fatalf("expected vcs build info in binary, got:\n%s", s)
	}
}
```
