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
    "strings"
    "testing"
    "time"

    "github.com/xhd2015/doctest/libdoc/testtree"
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
    testtree.WritePassFailTree(t, tmp, passCount, failCount)
    return tmp
}

const (
    singleFailLogMarker = "SINGLE_FAIL_LOG_MARKER"
    secondFailLogMarker = "SECOND_FAIL_LOG_MARKER"
    unwantedNonVerboseLogfMarker = "UNWANTED_NONVERBOSE_LOGF_MARKER"
)

func createRootTestTree(t *testing.T) string {
    t.Helper()
    tmp := t.TempDir()
    testtree.WriteFile(t, tmp, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))
    return tmp
}

func writePassLeaf(t *testing.T, root, name string) {
    t.Helper()
    leafDir := filepath.Join(root, name)
    os.MkdirAll(leafDir, 0755)
    os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(`import "testing"
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }`)), 0644)
    os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(`import "testing"
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}`)), 0644)
}

func writeFailLeafWithMarker(t *testing.T, root, name, marker string) {
    t.Helper()
    leafDir := filepath.Join(root, name)
    os.MkdirAll(leafDir, 0755)
    os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(`import "testing"
func Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }`)), 0644)
    assertCode := fmt.Sprintf(`import "testing"
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    t.Fatal(%q)
}`, marker)
    os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(assertCode)), 0644)
}

func createLogfPassTree(t *testing.T) string {
    t.Helper()
    tmp := createRootTestTree(t)
    leafDir := filepath.Join(tmp, "logf_leaf")
    os.MkdirAll(leafDir, 0755)
    setupCode := fmt.Sprintf(`import "testing"
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Logf(%q)
    return nil
}`, unwantedNonVerboseLogfMarker)
    os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(doctestGoBlock(setupCode)), 0644)
    os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(doctestGoBlock(`import "testing"
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}`)), 0644)
    return tmp
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

func findInlineSummaryLine(stdout string) string {
    for _, line := range strings.Split(stdout, "\n") {
        if strings.Contains(line, " Run, ") && strings.Contains(line, " Cached") {
            return line
        }
    }
    return ""
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
