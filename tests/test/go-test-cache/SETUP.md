# Scenario

**Feature**: the doctest binary is built by the root Setup

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The doctest binary is built by the root Setup.
- Tests run with an extended timeout to allow for Go compilation (first run is slow).

## Steps
1. Set a generous timeout (120s) for tests that compile Go code.
2. Provide shared helpers for creating temporary doc-style test projects.

```go
import (
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

var bt = "`" + "`" + "`"

func doctestGoBlock(code string) string {
    return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

func doctestBody(extraRunCode string) string {
    return "import \"testing\"\n\ntype Request struct{ Args []string; WorkDir string }\ntype Response struct{ ExitCode int; Stdout string; Stderr string }\n\n" + extraRunCode
}
func rootSetupContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}

func leafSetupContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}

func leafAssertContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}")
}

func createTestTree(dir string, extraRunCode string) error {
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody(extraRunCode))), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(rootSetupContent()), 0644); err != nil {
        return err
    }
    leafDir := filepath.Join(dir, "simple")
    if err := os.MkdirAll(leafDir, 0755); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(leafSetupContent()), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
        return err
    }
    return nil
}

func createTempTestProject(t *testing.T, dirName string) string {
    t.Helper()
    tmp := t.TempDir()
    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }
    testDir := filepath.Join(tmp, dirName)
    if err := os.MkdirAll(testDir, 0755); err != nil {
        t.Fatalf("mkdir test dir: %v", err)
    }
    if err := createTestTree(testDir, "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"); err != nil {
        t.Fatalf("create test tree: %v", err)
    }
    return testDir
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
