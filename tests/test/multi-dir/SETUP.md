# Scenario

**Feature**: a temporary project with `go.mod` and multiple DOCTEST.md test trees

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temporary project with `go.mod` and multiple DOCTEST.md test trees.
- The doctest binary is built by the root Setup.

## Steps
1. Create a temp project with several doctest test trees (`test_a`, `test_b`, `test_c`) and a non-test directory (`no_tests`).
2. Each leaf configures which directories to pass to the doctest binary.

```go
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

var bt = "`" + "`" + "`"

func doctestGoBlock(code string) string {
    return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

func doctestBody() string {
    return strings.Join([]string{
        "import \"testing\"",
        "",
        "type Request struct{ Args []string; Env []string; WorkDir string }",
        "type Response struct{ ExitCode int; Stdout string; Stderr string }",
        "",
        "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }",
    }, "\n")
}
func rootSetupContent() string {
    return "# Scenario\n\n**Feature**: multi-dir test tree root\n\n## Steps\n1. Root setup.\n\n" + doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error {\n    t.Logf(\"setup\")\n    return nil\n}")
}

func leafSetupContent() string {
    return "# Scenario\n\n**Feature**: multi-dir leaf\n\n## Steps\n1. Leaf setup.\n\n" + doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error {\n    t.Logf(\"setup\")\n    return nil\n}")
}

func leafAssertContent() string {
    return "## Expected\n- pass\n\n" + doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}")
}

func createTestTree(parent string, name string) error {
    root := filepath.Join(parent, name)
    if err := os.MkdirAll(filepath.Join(root, "simple"), 0755); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(root, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody())), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(root, "SETUP.md"), []byte(rootSetupContent()), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(root, "simple", "SETUP.md"), []byte(leafSetupContent()), 0644); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(root, "simple", "ASSERT.md"), []byte(leafAssertContent()), 0644); err != nil {
        return err
    }
    return nil
}

func createMultiDirProject(t *testing.T, req *Request) string {
    t.Helper()
    tmp := t.TempDir()

    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }

    if err := createTestTree(tmp, "test_a"); err != nil {
        t.Fatalf("create test_a: %v", err)
    }
    if err := createTestTree(tmp, "test_b"); err != nil {
        t.Fatalf("create test_b: %v", err)
    }
    if err := createTestTree(tmp, "test_c"); err != nil {
        t.Fatalf("create test_c: %v", err)
    }

    if err := os.MkdirAll(filepath.Join(tmp, "no_tests"), 0755); err != nil {
        t.Fatalf("mkdir no_tests: %v", err)
    }

    return tmp
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
