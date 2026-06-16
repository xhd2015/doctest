# Scenario

**Feature**: a doctest tree with a SETUP.md containing invalid Go code (unclosed string literal)

```
# Go import processing during code generation
doctest build -> parse imports -> remove unused -> report syntax errors
```

## Preconditions
- A doctest tree with a SETUP.md containing invalid Go code (unclosed string literal).
- `imports.Process` should fail and report a clean error.

## Steps
1. Create a temp Go project with `go.mod`.
2. Create a doctest root with a SETUP.md that has a syntax error (unclosed string `"unclosed).
3. Create a leaf with SETUP.md and ASSERT.md.
4. Set `req.Args` to run `doctest test -v <test-dir>`.
5. The Run() from root executes the binary; ASSERT.md checks for error.

```go
import (
    "os"
    "path/filepath"
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

    brokenCode := "func Run(t *testing.T, req *Request) (*Response, error) {\n    _ = \"unclosed\n    return &Response{}, nil\n}\n"
    rootSetup := rootSetupContent(brokenCode)

    if err := createDoctestRoot(testDir, rootSetup); err != nil {
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
