# Scenario

**Feature**: `build.Test` stderr paths are shortened via `DisplayPath`

```
# build.Test pipeline
build.Test(dir, opts) -> announceRoots -> doctest header -> cd preview -> go test

# display-only formatting at stderr call sites
genRoot/runDir/dir -> DisplayPath -> shortened stderr line
```

## Preconditions

- The `build` package is importable (`github.com/xhd2015/doctest/libdoc/build`).
- Each leaf creates a minimal temp doctest tree and captures stderr from `build.Test`.
- Backtick characters in embedded Go strings use `\x60` to avoid conflicting with
  the outer markdown code fence.

## Steps

1. Create a temp Go project with `go.mod` and a minimal doctest tree.
2. Optionally `chdir` to the project root so cwd-relative display applies.
3. Call `build.Test` with stderr captured in a buffer.
4. Parse `→`, `doctest:`, and `cd` lines into `Response`.

## Context

- Leaves differ only in `GenDir` / gen-dir source configuration.
- `RemoveTemp: true` avoids leaving generated trees on disk.

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
)

type Request struct {
	GenDir string
}

type Response struct {
	Stderr    string
	ArrowLine string
	HeaderLine string
	CdLine    string
	TestErr   error
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

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func createMinimalTree(t *testing.T, projRoot string) string {
	t.Helper()
	if err := os.WriteFile(filepath.Join(projRoot, "go.mod"), []byte("module shortdisplay\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	testRoot := filepath.Join(projRoot, "tests", "feature")
	writeFile(t, testRoot, "SETUP.md", ""+
		"## Preconditions\n- Minimal tree.\n\n"+
		"## Steps\n1. Run returns immediately.\n\n"+
		"\x60\x60\x60go\n"+
		"import \"testing\"\n"+
		"type Request struct{}\n"+
		"type Response struct{}\n"+
		"func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"+
		"\x60\x60\x60\n")
	writeFile(t, testRoot, "leaf/SETUP.md", ""+
		"## Steps\n1. No setup.\n\n"+
		"\x60\x60\x60go\n"+
		"import \"testing\"\n"+
		"func Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+
		"\x60\x60\x60\n")
	writeFile(t, testRoot, "leaf/ASSERT.md", ""+
		"## Expected\n- Passes.\n\n"+
		"\x60\x60\x60go\n"+
		"func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"+
		"\x60\x60\x60\n")
	return testRoot
}

func parseStderrLines(stderr string) (string, string, string) {
	var arrow, header, cd string
	for _, line := range strings.Split(stderr, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "→ "):
			arrow = trimmed
		case strings.HasPrefix(trimmed, "doctest: "):
			header = trimmed
		case strings.HasPrefix(trimmed, "cd "):
			cd = trimmed
		}
	}
	return arrow, header, cd
}

func Run(t *testing.T, req *Request) (*Response, error) {
	saveAndRestoreCwd(t)
	projRoot := t.TempDir()
	testRoot := createMinimalTree(t, projRoot)
	if err := os.Chdir(projRoot); err != nil {
		t.Fatal(err)
	}

	var stderr bytes.Buffer
	opts := core.Options{
		GenDir:     req.GenDir,
		RemoveTemp: true,
		Stderr:     &stderr,
	}
	testErr := build.Test(testRoot, opts)
	out := stderr.String()
	arrow, header, cd := parseStderrLines(out)

	return &Response{
		Stderr:     out,
		ArrowLine:  arrow,
		HeaderLine: header,
		CdLine:     cd,
		TestErr:    testErr,
	}, nil
}
```