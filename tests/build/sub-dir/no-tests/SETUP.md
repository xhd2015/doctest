## Preconditions
- A doc-style test tree exists.
- A sub-directory has no ASSERT.md files.

## Steps
1. Create a test tree with a no-leaf-dir that has a SETUP.md but no ASSERT.md descendants.
2. Run `doctest build <treeRoot>/no-leaf-dir`.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    treeRoot := t.TempDir()
    bt := string(rune(96))
    d := bt + bt + bt

    os.WriteFile(filepath.Join(treeRoot, "DOCTEST.md"), []byte("# sub-dir test tree\n"), 0644)
    os.WriteFile(filepath.Join(treeRoot, "SETUP.md"), []byte(
        d+"go\n"+
        "type Request struct{}\n"+
        "type Response struct{}\n"+
        "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"+
        d+"\n"), 0644)

    noLeafDir := filepath.Join(treeRoot, "no-leaf-dir")
    os.MkdirAll(noLeafDir, 0755)
    os.WriteFile(filepath.Join(noLeafDir, "SETUP.md"), []byte(
        d+"go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+d+"\n"), 0644)

    req.Args = []string{"build", noLeafDir}
    return nil
}
```
