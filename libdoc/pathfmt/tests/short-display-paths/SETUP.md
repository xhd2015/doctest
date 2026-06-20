# Scenario

**Feature**: `DisplayPath` shortens filesystem paths for display-only CLI output

```
# formatter pipeline
caller path string -> DisplayPath -> Abs normalize -> cwd/home rules -> display string

# cwd rules (checked first)
path == cwd -> "." | strict child of cwd -> "./" + rel

# home shorten (when not under cwd)
path under home -> "~" + suffix | otherwise -> absolute unchanged
```

## Preconditions

- The `pathfmt` package is importable (`github.com/xhd2015/doctest/libdoc/pathfmt`).
- `DisplayPath` is display-only; tests call it directly (not via CLI).
- Ancestor `Setup` functions may `chdir` to control cwd; root saves and restores
  the original working directory.

## Steps

1. Ancestor `Setup` chains configure `req.Path` and optionally change cwd.
2. Root `Run` calls `pathfmt.DisplayPath(req.Path)` and records cwd for assertions.

## Context

- Platform-native `filepath` separators are used in expectations.
- No symlink canonicalization is required.
- Grouping nodes under `cwd-relative` and `home-shorten` create temp dirs and
  chdir before leaves set concrete paths.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/pathfmt"
)

type Request struct {
	Path string
}

type Response struct {
	Display string
	Cwd     string
}

func saveAndRestoreCwd(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})
	return wd
}

func chdirTo(t *testing.T, dir string) {
	t.Helper()
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(abs); err != nil {
		t.Fatal(err)
	}
}

func mkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
}

func Run(t *testing.T, req *Request) (*Response, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	display := pathfmt.DisplayPath(req.Path)
	return &Response{Display: display, Cwd: cwd}, nil
}
```