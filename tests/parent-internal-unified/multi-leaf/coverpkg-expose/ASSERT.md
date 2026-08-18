## Expected

**Desired product behavior** (GREEN after profile-strip fix; **RED** on current code):

- `RunErr` empty: nested subject tree completes under cover + coverpkg.
- Cover profile at `req.CoverPath` exists and is non-empty.
- Combined stdout/stderr do **not** contain the classic expose cover failure
  during the suite run: `cover:` + `__doctest_internal_expose` + `no such file or directory`.
- Prefer single suite package args when the go test line is present.
- **Final coverprofile must not list session-generated expose facades** under
  `modPath/__doctest_internal_expose/` (only real product packages such as
  `internal/greet`).
- From `ModuleRoot`, `go tool cover -func=<CoverPath>` exits 0 (scaff CI
  report step after merge).

Today (post materialize-for-cover): suite PASS and profile non-empty, but the
profile still contains `…/__doctest_internal_expose/…/expose.go` lines; after
cleanup, `go tool cover -func` fails with
`no required module provides package …/__doctest_internal_expose/…` → **RED**.

```go
import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	_ = d
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if req.CoverPath == "" {
		t.Fatal("CoverPath empty")
	}
	if req.CoverPkg == "" {
		t.Fatal("CoverPkg empty — this leaf requires product-module coverpkg")
	}
	if req.ModuleRoot == "" {
		t.Fatal("ModuleRoot empty — need product cwd for go tool cover")
	}

	combined := resp.Stdout + "\n" + resp.Stderr + "\n" + resp.RunErr
	exposeCoverOpen := strings.Contains(combined, "__doctest_internal_expose") &&
		strings.Contains(combined, "no such file or directory") &&
		(strings.Contains(combined, "cover:") || strings.Contains(combined, "expose.go"))

	if resp.RunErr != "" {
		hint := ""
		if exposeCoverOpen {
			hint = " (expose: go tool cover opens overlay-only expose.go)"
		}
		t.Fatalf("want cover+coverpkg success for parent-internal expose, got RunErr=%s%s\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, hint, resp.Stdout, resp.Stderr)
	}
	if exposeCoverOpen {
		t.Fatalf("cover must not open missing __doctest_internal_expose.go under coverpkg=%q\nstdout:\n%s\nstderr:\n%s",
			req.CoverPkg, resp.Stdout, resp.Stderr)
	}

	info, statErr := os.Stat(req.CoverPath)
	if statErr != nil {
		t.Fatalf("expected coverprofile at %s: %v\nstdout:\n%s\nstderr:\n%s",
			req.CoverPath, statErr, resp.Stdout, resp.Stderr)
	}
	if info.Size() == 0 {
		t.Fatalf("coverprofile %s is empty", req.CoverPath)
	}
	if len(resp.GoTestPackageArgs) > 1 {
		t.Fatalf("coverpkg path should be single-package suite run, got pkgs=%v line=%q",
			resp.GoTestPackageArgs, resp.GoTestDisplayLine)
	}
	if len(resp.GoTestPackageArgs) == 1 && !strings.Contains(resp.GoTestPackageArgs[0], "suite") {
		t.Fatalf("single package should be suite, got %q", resp.GoTestPackageArgs[0])
	}

	// Desired: final profile excludes session-generated expose facades for this
	// product module (crime scene A1 inverted). Prefix is the known expose import
	// root for modPath — not a blind global string match on unrelated modules.
	exposeFilePrefix := modPath + "/" + "__doctest_internal_expose" + "/"
	profile, readErr := os.ReadFile(req.CoverPath)
	if readErr != nil {
		t.Fatalf("read coverprofile %s: %v", req.CoverPath, readErr)
	}
	var exposeLines []string
	for _, line := range strings.Split(string(profile), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "mode:") {
			continue
		}
		if strings.HasPrefix(line, exposeFilePrefix) {
			exposeLines = append(exposeLines, line)
		}
	}
	if len(exposeLines) > 0 {
		t.Fatalf("coverprofile must not list session-generated expose facades (prefix %q); got %d line(s):\n%s\nfull profile:\n%s",
			exposeFilePrefix, len(exposeLines), strings.Join(exposeLines, "\n"), string(profile))
	}

	// Desired: scaff-style report step works after doctest returns (crime scene A3).
	cmd := exec.Command("go", "tool", "cover", "-func", req.CoverPath)
	cmd.Dir = req.ModuleRoot
	out, coverErr := cmd.CombinedOutput()
	if coverErr != nil {
		t.Fatalf("go tool cover -func=%s (cwd=%s) must succeed after doctest; err=%v\noutput:\n%s\nprofile:\n%s",
			req.CoverPath, req.ModuleRoot, coverErr, string(out), string(profile))
	}
}
```
