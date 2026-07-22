# Scenario

**Feature**: the root Setup has built the doctest binary and set req.Bin

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The root Setup has built the doctest binary and set req.Bin.
- Tests need extended timeout for Go compilation (two doctest invocations).
- **Parallel-safe**: multi-run results live on `req.MRFirst` / `req.MRSecond`
  (not package-global `isoState`).

## Steps
1. Provide a multi-run `doMultiRun` that creates a test tree with groups and runs
   doctest twice for a given scenario string.
2. Each leaf passes its scenario to `doMultiRun`; results are written onto `req`.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

var bt = "`" + "`" + "`"

func doctestGoBlock(code string) string {
	return bt + "go\n" + code + bt + "\n"
}

func createMultiGroupTree(t *testing.T, dir string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(testtree.MinimalRunGo())), 0644); err != nil {
		t.Fatalf("write root DOCTEST.md: %v", err)
	}

	groups := []struct {
		name  string
		leaves []string
	}{
		{"group-a", []string{"leaf-1", "leaf-2"}},
		{"group-b", []string{"leaf-3"}},
	}

	for _, g := range groups {
		gDir := filepath.Join(dir, g.name)
		if err := os.MkdirAll(gDir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", g.name, err)
		}
		groupSetupCode := "import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"
		if err := os.WriteFile(filepath.Join(gDir, "SETUP.md"), []byte(doctestGoBlock(groupSetupCode)), 0644); err != nil {
			t.Fatalf("write %s SETUP.md: %v", g.name, err)
		}

		for _, leaf := range g.leaves {
			lDir := filepath.Join(gDir, leaf)
			if err := os.MkdirAll(lDir, 0755); err != nil {
				t.Fatalf("mkdir %s/%s: %v", g.name, leaf, err)
			}
			leafSetupCode := "import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"
			if err := os.WriteFile(filepath.Join(lDir, "SETUP.md"), []byte(doctestGoBlock(leafSetupCode)), 0644); err != nil {
				t.Fatalf("write %s/%s SETUP.md: %v", g.name, leaf, err)
			}
			leafAssertCode := "import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"
			if err := os.WriteFile(filepath.Join(lDir, "ASSERT.md"), []byte(doctestGoBlock(leafAssertCode)), 0644); err != nil {
				t.Fatalf("write %s/%s ASSERT.md: %v", g.name, leaf, err)
			}
		}
	}
}

func doRun(t *testing.T, bin string, args []string, timeout time.Duration) *Response {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Env = os.Environ()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	resp := &Response{}
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
		} else {
			resp.Err = err
		}
	}
	resp.Stdout = stdoutBuf.String()
	resp.Stderr = stderrBuf.String()
	return resp
}

func doMultiRun(t *testing.T, req *Request, scenario string) {
    if req.Bin == "" {
        t.Fatalf("req.Bin is not set")
    }

    treeRoot := t.TempDir()
    createMultiGroupTree(t, treeRoot)

    var firstArgs, secondArgs []string
    switch scenario {
    case "subdir_after_full":
        firstArgs = []string{"test", "-v", treeRoot}
        secondArgs = []string{"test", "-v", filepath.Join(treeRoot, "group-a")}
    case "full_after_subdir":
        firstArgs = []string{"test", "-v", filepath.Join(treeRoot, "group-a")}
        secondArgs = []string{"test", "-v", treeRoot}
    case "subdir_after_subdir":
        firstArgs = []string{"test", "-v", filepath.Join(treeRoot, "group-a")}
        secondArgs = []string{"test", "-v", filepath.Join(treeRoot, "group-b")}
    default:
        t.Fatalf("unknown scenario: %s", scenario)
    }

    runTimeout := 120 * time.Second
    req.MRFirst = doRun(t, req.Bin, firstArgs, runTimeout)
    req.MRSecond = doRun(t, req.Bin, secondArgs, runTimeout)
}

func Setup(t *testing.T, req *Request) error {
    _ = fmt.Sprintf
    req.MRFirst = nil
    req.MRSecond = nil
    req.Timeout = 150 * time.Second
    req.Args = []string{}
    return nil
}
```
