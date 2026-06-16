# Scenario

**Feature**: a doctest tree with a SETUP.md that imports `"fmt"` but never calls any `fmt.*` function

```
# Go import processing during code generation
doctest build -> parse imports -> remove unused -> report syntax errors
```

## Preconditions
- A doctest tree with a SETUP.md that imports `"fmt"` but never calls any `fmt.*` function.
- `imports.Process` should remove the unused import from generated Go code.

## Steps
1. Create a temp Go project with `go.mod`.
2. Create a doctest root with a SETUP.md that imports `"fmt"` (unused) and `"testing"`.
3. Create a leaf with SETUP.md and ASSERT.md that don't use `"fmt"`.
4. Set `req.Args` to run `doctest test -v <test-dir>`.
5. The Run() from root executes the binary; ASSERT.md checks exit code.

```go
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    proj := t.TempDir()
    if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }

    testDir := filepath.Join(proj, "tests")
    if err := os.MkdirAll(testDir, 0755); err != nil {
        t.Fatalf("mkdir test dir: %v", err)
    }

    runCode := "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
    rootSetup := rootSetupContent(runCode)
    rootSetupLines := strings.Split(rootSetup, "\n")
    var newLines []string
    for _, line := range rootSetupLines {
        newLines = append(newLines, line)
        if line == "import \"testing\"" {
            newLines[len(newLines)-1] = "import (\n    \"testing\"\n    \"fmt\"\n)"
        }
    }
    rootSetupWithUnusedImport := strings.Join(newLines, "\n")

    if err := createDoctestRoot(testDir, rootSetupWithUnusedImport); err != nil {
        t.Fatalf("create doctest root: %v", err)
    }

    leafDir := filepath.Join(testDir, "leaf")
    if err := createDoctestLeaf(leafDir, leafSetupContent()); err != nil {
        t.Fatalf("create leaf: %v", err)
    }

    req.Args = []string{"test", "-v", testDir}
    return nil
}
```
