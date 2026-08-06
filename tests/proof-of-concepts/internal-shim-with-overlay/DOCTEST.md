# Go internal-package shim + overlay (proof of concept)

## Version

0.0.1

## Layer

**L3 e2e** — runs the host `go` toolchain against a fixture module under
`testdata/realmod/`. Leaves carry `label: e2e` (skipped by default
discovery; use `--label e2e` or `--label-all`).

## Why this tree exists

Doctest unified gen blank-imports leaf packages from `__allleaves`. A
scenario path segment named `internal` becomes a Go **internal package**
import path, so importers outside the parent of `internal` fail with:

```text
use of internal package …/http/internal/… not allowed
```

This tree locks a **language-level** fix shape (not doctest gen yet):

1. Same-module **bridge package** under the parent of `internal`
   (`…/http/__doctest_internal_shim/…`) may import the internal leaf.
2. Code outside `http/` imports only the bridge.
3. The bridge may exist **on disk** or only via **`go -overlay`** (virtual file).
4. A bridge placed **outside** the parent of `internal` still fails.

## DSN

### Participants

- **Fixture module** — `testdata/realmod` (`example.com/realmod`), copied per
  leaf into `t.TempDir()`.
- **Internal leaf** — `http/internal/leaf` (restricted).
- **On-disk shim** — `http/__doctest_internal_shim/leaf` re-exports `Hello`.
- **Overlay shim** — virtual path
  `http/__doctest_internal_shim_overlay/leaf` mapped via `overlay.json`.
- **Wrong shim** — virtual `__wrong_shim/leaf` outside `http/` (negative).
- **`go` CLI** — `build` / `run` / `list` with optional `-overlay`.

### Behaviors

| Scenario | Action | Expect |
|----------|--------|--------|
| direct | `go build ./suite_direct` | fail: internal not allowed |
| on-disk-shim | `go run ./suite` | print `from-internal-leaf` |
| overlay-shim | `go run -overlay=overlay.json ./suite_overlay` | print `from-internal-leaf` |
| wrong-shim | `go build -overlay=… ./suite_wrong` | fail: internal not allowed |
| overlay-list | `go list -overlay=… <virtual pkg>` | list path, exit 0 |

## Decision tree

```text
internal-shim-with-overlay/
├── testdata/realmod/           # fixture (skipped by discovery)
├── direct-import-fails/
├── on-disk-shim-passes/
├── overlay-shim-passes/
├── wrong-shim-fails/
└── overlay-list-ok/
```

## How to run

```sh
# skipped by default (e2e)
doctest test ./tests/proof-of-concepts/internal-shim-with-overlay

# run proof
doctest test --label e2e -v ./tests/proof-of-concepts/internal-shim-with-overlay/...
# or
doctest test --label-all -v ./tests/proof-of-concepts/internal-shim-with-overlay/...
```

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

// Request is filled by root + leaf Setup.
type Request struct {
	// WorkDir is a per-leaf copy of testdata/realmod (absolute).
	WorkDir string
	// Scenario selects the go invocation (set by leaf Setup).
	Scenario string
}

// Response captures go CLI outcome.
type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Combined string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req == nil || req.WorkDir == "" {
		return nil, fmt.Errorf("WorkDir is required")
	}
	if req.Scenario == "" {
		return nil, fmt.Errorf("Scenario is required")
	}

	outBin := filepath.Join(req.WorkDir, "out.bin")
	// -overlay path must be absolute: when cmd.Dir is set, a relative overlay
	// path is resolved against the process cwd, not WorkDir.
	overlayFlag := "-overlay=" + filepath.Join(req.WorkDir, "overlay.json")
	var args []string
	switch req.Scenario {
	case "direct":
		args = []string{"build", "-o", outBin, "./suite_direct"}
	case "on-disk-shim":
		args = []string{"run", "./suite"}
	case "overlay-shim":
		args = []string{"run", overlayFlag, "./suite_overlay"}
	case "wrong-shim":
		args = []string{"build", overlayFlag, "-o", outBin, "./suite_wrong"}
	case "overlay-list":
		args = []string{
			"list",
			overlayFlag,
			"example.com/realmod/http/__doctest_internal_shim_overlay/leaf",
		}
	default:
		return nil, fmt.Errorf("unknown Scenario %q", req.Scenario)
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	resp := &Response{
		Combined: string(out),
		Stdout:   string(out),
		Stderr:   string(out),
	}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			resp.ExitCode = ee.ExitCode()
		} else {
			resp.ExitCode = 1
		}
		// Expected toolchain failures are asserted on ExitCode/Combined.
		return resp, nil
	}
	resp.ExitCode = 0
	return resp, nil
}

// materializeFixture copies testdata/realmod into dest and writes overlay.json
// with absolute Replace keys for this workdir.
func materializeFixture(t *testing.T, fixtureRoot, dest string) {
	t.Helper()
	if err := copyDir(fixtureRoot, dest); err != nil {
		t.Fatalf("copy fixture: %v", err)
	}
	writeOverlayJSON(t, dest)
}

func writeOverlayJSON(t *testing.T, workDir string) {
	t.Helper()
	abs, err := filepath.Abs(workDir)
	if err != nil {
		t.Fatalf("abs workdir: %v", err)
	}
	// Prefer the path form `go list -m` uses for this module dir. On macOS,
	// filepath.Abs often yields /private/var/... while the go command uses
	// /var/...; -overlay Replace keys must match the go command's view.
	modDir := goModuleDir(t, workDir)
	if modDir == "" {
		modDir = stripPrivatePrefix(abs)
	}
	pairs := [][2]string{
		{
			filepath.Join(modDir, "http", "__doctest_internal_shim_overlay", "leaf", "shim.go"),
			filepath.Join(modDir, "overlay-src", "overlay_shim.go"),
		},
		{
			filepath.Join(modDir, "__wrong_shim", "leaf", "shim.go"),
			filepath.Join(modDir, "overlay-src", "wrong_shim.go"),
		},
	}
	// Also register /private and non-/private variants so either view matches.
	rep := map[string]string{}
	for _, pair := range pairs {
		for _, k := range darwinPathVariants(pair[0]) {
			// Backing file: prefer a form that exists on disk.
			back := firstExisting(darwinPathVariants(pair[1])...)
			if back == "" {
				back = pair[1]
			}
			rep[k] = back
		}
	}
	body, err := json.MarshalIndent(map[string]map[string]string{"Replace": rep}, "", "  ")
	if err != nil {
		t.Fatalf("marshal overlay: %v", err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(filepath.Join(workDir, "overlay.json"), body, 0644); err != nil {
		t.Fatalf("write overlay.json: %v", err)
	}
}

// goModuleDir returns `go list -m -f {{.Dir}}` for workDir, or "".
func goModuleDir(t *testing.T, workDir string) string {
	t.Helper()
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	cmd.Dir = workDir
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func stripPrivatePrefix(p string) string {
	if strings.HasPrefix(p, "/private/") {
		return strings.TrimPrefix(p, "/private")
	}
	return p
}

func darwinPathVariants(p string) []string {
	p = filepath.Clean(p)
	seen := map[string]bool{}
	var out []string
	add := func(s string) {
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		out = append(out, s)
	}
	add(p)
	if strings.HasPrefix(p, "/private/") {
		add(strings.TrimPrefix(p, "/private"))
	} else if strings.HasPrefix(p, "/var/") || strings.HasPrefix(p, "/tmp/") {
		add("/private" + p)
	}
	return out
}

func firstExisting(paths ...string) string {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(dst, 0755)
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}

func containsInternalDenied(s string) bool {
	return strings.Contains(s, "use of internal package") ||
		strings.Contains(s, "internal package") && strings.Contains(s, "not allowed")
}
```
