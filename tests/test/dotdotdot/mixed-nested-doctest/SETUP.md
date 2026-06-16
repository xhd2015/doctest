# Scenario

**Feature**: a temporary project with a complex doctest tree structure

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A temporary project with a complex doctest tree structure.
- `ancestor/` has DOCTEST.md (the parent doctest root).
- `ancestor/leaf/` has ASSERT.md (a test leaf within the ancestor tree).
- `ancestor/leaf/nested-sub2/` has DOCTEST.md (a nested doctest root).

```go
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

func mkGoBlock(code string) string {
    return "## Test\n\n" + bt + "go\n" + code + "\n" + bt + "\n"
}

func createMixedTestProject(t *testing.T) string {
    t.Helper()
    projDir := t.TempDir()

    os.WriteFile(filepath.Join(projDir, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644)

    // ancestor DOCTEST root with non-stub Setup
    ancestorDir := filepath.Join(projDir, "ancestor")
    os.MkdirAll(ancestorDir, 0755)
    os.WriteFile(filepath.Join(ancestorDir, "DOCTEST.md"), []byte("# Ancestor\n"), 0644)
    os.WriteFile(filepath.Join(ancestorDir, "SETUP.md"), []byte(mkGoBlock(strings.Join([]string{
        "import \"testing\"",
        "type Request struct{}",
        "type Response struct{}",
        "func Setup(t *testing.T, req *Request) error { t.Logf(\"ancestor setup\"); return nil }",
        "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }",
    }, "\n"))), 0644)

    // leaf directory with ASSERT.md (test case within ancestor tree)
    leafDir := filepath.Join(ancestorDir, "leaf")
    os.MkdirAll(leafDir, 0755)
    os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(mkGoBlock(
        "import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { t.Logf(\"leaf setup\"); return nil }",
    )), 0644)
    os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(mkGoBlock(
        "import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}",
    )), 0644)

    // nested DOCTEST root: leaf/nested-sub2 (non-stub)
    sub2Dir := filepath.Join(leafDir, "nested-sub2")
    os.MkdirAll(sub2Dir, 0755)
    os.WriteFile(filepath.Join(sub2Dir, "DOCTEST.md"), []byte("# Nested Sub2\n"), 0644)
    os.WriteFile(filepath.Join(sub2Dir, "SETUP.md"), []byte(mkGoBlock(strings.Join([]string{
        "import \"testing\"",
        "type Request struct{}",
        "type Response struct{}",
        "func Setup(t *testing.T, req *Request) error { t.Logf(\"sub2 setup\"); return nil }",
        "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }",
    }, "\n"))), 0644)

    // nested-sub2's test leaf
    simpleDir := filepath.Join(sub2Dir, "simple")
    os.MkdirAll(simpleDir, 0755)
    os.WriteFile(filepath.Join(simpleDir, "SETUP.md"), []byte(mkGoBlock(
        "import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { t.Logf(\"sub2 simple setup\"); return nil }",
    )), 0644)
    os.WriteFile(filepath.Join(simpleDir, "ASSERT.md"), []byte(mkGoBlock(
        "import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}",
    )), 0644)

    return projDir
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
