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

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
