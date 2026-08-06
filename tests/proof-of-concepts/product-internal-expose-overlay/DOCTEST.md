# Product internal expose via overlay (kind B proof of concept)

## Version

0.0.1

## Layer

**L3 e2e** — host `go` only (no doctest gen engine). Leaves: `label: e2e`.

## Why this tree exists

**Kind B:** a consumer module (`example.com/runner`, stand-in for mapping-gen
`testcase`) must not import `example.com/app/internal/…`. Fix shape:

1. Virtual package under the **product** module path, **without** the segment
   `internal`, e.g. `example.com/app/__doctest_internal_expose/greet`.
2. That package imports product `internal` and re-exports what tests need.
3. Materialize the package only via **`go -overlay`** so the product tree stays
   clean on disk.
4. Consumer imports the expose path (legal from outside).

Sibling of `internal-shim-with-overlay` (kind A, same-module gen shim).

## DSN

### Participants

- **app** — `example.com/app` with `internal/greet`.
- **expose body** — `expose-src/expose.go` (backing file for overlay).
- **virtual path** — `<appRoot>/__doctest_internal_expose/greet/expose.go` (not on disk).
- **runner** — `example.com/runner` with `replace example.com/app => ../app`.
- **suite_direct** — imports product internal → must fail.
- **suite_expose** — imports expose path with `-overlay` → must pass.

## Decision tree

```text
product-internal-expose-overlay/
├── testdata/workspace/
│   ├── app/internal/greet/
│   ├── expose-src/
│   └── runner/{suite_direct,suite_expose}/
├── direct-product-internal-fails/   # e2e
└── expose-via-overlay-passes/       # e2e
```

## How to run

```sh
go run ./cmd/doctest test --label e2e -count=1 -v \
  ./tests/proof-of-concepts/product-internal-expose-overlay/...
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

type Request struct {
	// WorkDir is the runner module root (absolute).
	WorkDir string
	// AppRoot is the product module root (absolute).
	AppRoot string
	// OverlayPath is absolute path to overlay.json under the work copy.
	OverlayPath string
	Scenario    string
}

type Response struct {
	ExitCode int
	Combined string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req == nil || req.WorkDir == "" {
		return nil, fmt.Errorf("WorkDir required")
	}
	if req.Scenario == "" {
		return nil, fmt.Errorf("Scenario required")
	}

	var args []string
	switch req.Scenario {
	case "direct":
		// Consumer → product internal (illegal).
		args = []string{"build", "-o", filepath.Join(req.WorkDir, "out.bin"), "./suite_direct"}
	case "expose-overlay":
		if req.OverlayPath == "" {
			return nil, fmt.Errorf("OverlayPath required")
		}
		args = []string{"run", "-overlay=" + req.OverlayPath, "./suite_expose"}
	default:
		return nil, fmt.Errorf("unknown Scenario %q", req.Scenario)
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	resp := &Response{Combined: string(out)}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			resp.ExitCode = ee.ExitCode()
		} else {
			resp.ExitCode = 1
		}
		return resp, nil
	}
	resp.ExitCode = 0
	return resp, nil
}

func materializeWorkspace(t *testing.T, fixtureRoot, dest string) {
	t.Helper()
	if err := copyDir(fixtureRoot, dest); err != nil {
		t.Fatalf("copy workspace: %v", err)
	}
	// Overlay first: suite_expose imports a virtual product package; tidy/list
	// need -overlay to see example.com/app/__doctest_internal_expose/greet.
	writeProductExposeOverlay(t, dest)
	runner := filepath.Join(dest, "runner")
	overlay := filepath.Join(dest, "overlay.json")
	tidy := exec.Command("go", "mod", "tidy", "-overlay="+overlay)
	tidy.Dir = runner
	tidy.Env = append(os.Environ(), "GOFLAGS=")
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy runner: %v\n%s", err, out)
	}
}

func writeProductExposeOverlay(t *testing.T, workspace string) {
	t.Helper()
	appRoot := filepath.Join(workspace, "app")
	// Prefer path form go list uses for the product module.
	modDir := goModuleDir(t, appRoot)
	if modDir == "" {
		modDir = stripPrivatePrefix(mustAbs(t, appRoot))
	}
	body := filepath.Join(workspace, "expose-src", "expose.go")
	bodyAbs := firstExisting(darwinPathVariants(mustAbs(t, body))...)
	if bodyAbs == "" {
		bodyAbs = mustAbs(t, body)
	}
	virtRel := filepath.Join("__doctest_internal_expose", "greet", "expose.go")
	rep := map[string]string{}
	for _, root := range darwinPathVariants(modDir) {
		rep[filepath.Join(root, virtRel)] = bodyAbs
	}
	// Also ensure body path variants exist as values when needed
	for k := range rep {
		if _, err := os.Stat(rep[k]); err != nil {
			// try other body forms
			if alt := firstExisting(darwinPathVariants(mustAbs(t, body))...); alt != "" {
				rep[k] = alt
			}
		}
	}
	raw, err := json.MarshalIndent(map[string]map[string]string{"Replace": rep}, "", "  ")
	if err != nil {
		t.Fatalf("marshal overlay: %v", err)
	}
	raw = append(raw, '\n')
	overlayPath := filepath.Join(workspace, "overlay.json")
	if err := os.WriteFile(overlayPath, raw, 0644); err != nil {
		t.Fatalf("write overlay: %v", err)
	}
}

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

func mustAbs(t *testing.T, p string) string {
	t.Helper()
	a, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("abs %s: %v", p, err)
	}
	return a
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
		(strings.Contains(s, "internal package") && strings.Contains(s, "not allowed"))
}
```
