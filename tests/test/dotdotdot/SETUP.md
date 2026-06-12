## Preconditions
- The doctest binary is built by the root Setup.

## Steps
1. Create a temporary project with `go.mod` and two `DOCTEST.md` trees.
2. The project also has a nested submodule with its own `DOCTEST.md` (must not be discovered).
3. Run `doctest test ./...` from the project root or a subdir.
4. Verify all test cases pass and nested module tests are excluded.

```go
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

var bt = "`" + "`" + "`"

func doctestGoBlock(code string) string {
    return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

func rootSetupContent() string {
    code := strings.Join([]string{
        "import \"testing\"",
        "",
        "type Request struct{ Args []string; Env []string; WorkDir string }",
        "type Response struct{ ExitCode int; Stdout string; Stderr string }",
        "",
        "func Setup(t *testing.T, req *Request) error {",
        "    t.Logf(\"setup\")",
        "    return nil",
        "}",
        "",
        "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }",
    }, "\n")
    return doctestGoBlock(code)
}

func leafSetupContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error {\n    t.Logf(\"setup\")\n    return nil\n}")
}

func leafAssertContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}")
}

func createTestTree(parent string, name string) error {
    root := filepath.Join(parent, name)
    if err := os.MkdirAll(filepath.Join(root, "simple"), 0755); err != nil {
        return err
    }
    if err := os.WriteFile(filepath.Join(root, "DOCTEST.md"), []byte("# "+name+" Tests\n"), 0644); err != nil {
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

func createTempProject(t *testing.T, req *Request) string {
    t.Helper()
    tmp := t.TempDir()

    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }

    if err := createTestTree(tmp, "alpha_test"); err != nil {
        t.Fatalf("create alpha_test: %v", err)
    }
    if err := createTestTree(tmp, "beta_test"); err != nil {
        t.Fatalf("create beta_test: %v", err)
    }

    nestedRoot := filepath.Join(tmp, "nested")
    if err := os.MkdirAll(nestedRoot, 0755); err != nil {
        t.Fatalf("mkdir nested: %v", err)
    }
    if err := os.WriteFile(filepath.Join(nestedRoot, "go.mod"), []byte("module nested\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write nested go.mod: %v", err)
    }
    if err := createTestTree(nestedRoot, "hidden_test"); err != nil {
        t.Fatalf("create hidden_test: %v", err)
    }

    return tmp
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
