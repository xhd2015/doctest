# Scenario

**Feature**: test trees are created programmatically as temp directories

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- Test trees are created programmatically as temp directories.
- Each leaf creates its own test tree and configures the doctest args.

## Steps
1. Write test tree files to a temp directory.
2. Run `doctest test` on the temp directory.
3. Verify dots and summary in stdout.

```go
import (
    "fmt"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func bt(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = '`'
    }
    return string(b)
}

func doctestGoBlock(code string) string {
    fence := bt(3)
    return "## Test\n\n" + fence + "go\n" + code + "\n" + fence + "\n"
}

func createPassFailTree(t *testing.T, passCount int, failCount int) string {
    t.Helper()
    tmp := t.TempDir()

    rootSetup := `import "testing"

type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }`
    os.WriteFile(filepath.Join(tmp, "SETUP.md"), []byte(doctestGoBlock(rootSetup)), 0644)
    os.WriteFile(filepath.Join(tmp, "DOCTEST.md"), []byte("# progress-dots test tree\n"), 0644)

    for i := 0; i < passCount; i++ {
        name := fmt.Sprintf("pass_%d", i+1)
        leafDir := filepath.Join(tmp, name)
        os.MkdirAll(leafDir, 0755)
        os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(`import "testing"
func Setup(t *testing.T, req *Request) error { _ = req; return nil }`)), 0644)
        os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(`import "testing"
func Assert(t *testing.T, req *Request, resp *Response, err error) {}`)), 0644)
    }

    for i := 0; i < failCount; i++ {
        name := fmt.Sprintf("fail_%d", i+1)
        leafDir := filepath.Join(tmp, name)
        os.MkdirAll(leafDir, 0755)
        os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(`import "testing"
func Setup(t *testing.T, req *Request) error { _ = req; return nil }`)), 0644)
        os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(`import "testing"
func Assert(t *testing.T, req *Request, resp *Response, err error) {
    t.Fatal("forced failure")
}`)), 0644)
    }

    return tmp
}

const (
    singleFailLogMarker = "SINGLE_FAIL_LOG_MARKER"
    secondFailLogMarker = "SECOND_FAIL_LOG_MARKER"
)

func createRootTestTree(t *testing.T) string {
    t.Helper()
    tmp := t.TempDir()
    rootSetup := `import "testing"

type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }`
    os.WriteFile(filepath.Join(tmp, "SETUP.md"), []byte(doctestGoBlock(rootSetup)), 0644)
    os.WriteFile(filepath.Join(tmp, "DOCTEST.md"), []byte("# progress-dots test tree\n"), 0644)
    return tmp
}

func writePassLeaf(t *testing.T, root, name string) {
    t.Helper()
    leafDir := filepath.Join(root, name)
    os.MkdirAll(leafDir, 0755)
    os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(`import "testing"
func Setup(t *testing.T, req *Request) error { _ = req; return nil }`)), 0644)
    os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(`import "testing"
func Assert(t *testing.T, req *Request, resp *Response, err error) {}`)), 0644)
}

func writeFailLeafWithMarker(t *testing.T, root, name, marker string) {
    t.Helper()
    leafDir := filepath.Join(root, name)
    os.MkdirAll(leafDir, 0755)
    os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(`import "testing"
func Setup(t *testing.T, req *Request) error { _ = req; return nil }`)), 0644)
    assertCode := fmt.Sprintf(`import "testing"
func Assert(t *testing.T, req *Request, resp *Response, err error) {
    t.Fatal(%q)
}`, marker)
    os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(assertCode)), 0644)
}

func createSingleFailLogTree(t *testing.T) string {
    t.Helper()
    tmp := createRootTestTree(t)
    writeFailLeafWithMarker(t, tmp, "fail_1", singleFailLogMarker)
    return tmp
}

func createSecondOfThreeFailsTree(t *testing.T) string {
    t.Helper()
    tmp := createRootTestTree(t)
    writePassLeaf(t, tmp, "pass_1")
    writeFailLeafWithMarker(t, tmp, "fail_2", secondFailLogMarker)
    writePassLeaf(t, tmp, "pass_3")
    return tmp
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
