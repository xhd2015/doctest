## Preconditions
- A doc-style test tree exists with multiple leaves.
- Build a specific leaf directory.

## Steps
1. Create a test tree with root DOCTEST.md, SETUP.md, leaves group-a/leaf-1 and group-a/leaf-2.
2. Run `doctest build <treeRoot>/group-a/leaf-1`.

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

    groupSetup := d+"go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+d+"\n"
    leafSetup := d+"go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+d+"\n"
    leafAssert := d+"go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"+d+"\n"

    ga := filepath.Join(treeRoot, "group-a")
    os.MkdirAll(ga, 0755)
    os.WriteFile(filepath.Join(ga, "SETUP.md"), []byte(groupSetup), 0644)

    leaf1 := filepath.Join(ga, "leaf-1")
    os.MkdirAll(leaf1, 0755)
    os.WriteFile(filepath.Join(leaf1, "SETUP.md"), []byte(leafSetup), 0644)
    os.WriteFile(filepath.Join(leaf1, "ASSERT.md"), []byte(leafAssert), 0644)

    leaf2 := filepath.Join(ga, "leaf-2")
    os.MkdirAll(leaf2, 0755)
    os.WriteFile(filepath.Join(leaf2, "SETUP.md"), []byte(leafSetup), 0644)
    os.WriteFile(filepath.Join(leaf2, "ASSERT.md"), []byte(leafAssert), 0644)

    req.Args = []string{"build", leaf1}
    return nil
}
```
