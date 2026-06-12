## Preconditions
- Two Go modules exist: mod-a has DOCTEST.md + SETUP.md, mod-b has its own SETUP.md.
- Running on a sub-dir of mod-b should not use mod-a's DOCTEST.md.

## Steps
1. Create mod-a (go.mod, tests/DOCTEST.md, tests/SETUP.md, leaf-a).
2. Create mod-b (go.mod, SETUP.md with types+Run, sub/leaf-b).
3. Run `doctest test <mod-b>/sub/leaf-b`.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    tmp := t.TempDir()
    bt := string(rune(96))
    d := bt + bt + bt

    modA := filepath.Join(tmp, "mod-a")
    modATests := filepath.Join(modA, "tests")
    os.MkdirAll(modATests, 0755)
    os.WriteFile(filepath.Join(modA, "go.mod"), []byte("module mod-a\n\ngo 1.21\n"), 0644)
    os.WriteFile(filepath.Join(modATests, "DOCTEST.md"), []byte("# mod-a tests\n"), 0644)
    os.WriteFile(filepath.Join(modATests, "SETUP.md"), []byte(
        d+"go\n"+
        "type RequestA struct{}\n"+
        "type ResponseA struct{}\n"+
        "func Run(t *testing.T, req *RequestA) (*ResponseA, error) { return &ResponseA{}, nil }\n"+
        d+"\n"), 0644)
    leafA := filepath.Join(modATests, "leaf-a")
    os.MkdirAll(leafA, 0755)
    os.WriteFile(filepath.Join(leafA, "SETUP.md"), []byte(d+"go\nfunc Setup(t *testing.T, req *RequestA) error { _ = req; return nil }\n"+d+"\n"), 0644)
    os.WriteFile(filepath.Join(leafA, "ASSERT.md"), []byte(d+"go\nfunc Assert(t *testing.T, req *RequestA, resp *ResponseA, err error) {}\n"+d+"\n"), 0644)

    modB := filepath.Join(tmp, "mod-b")
    os.MkdirAll(modB, 0755)
    os.WriteFile(filepath.Join(modB, "go.mod"), []byte("module mod-b\n\ngo 1.21\n"), 0644)
    os.WriteFile(filepath.Join(modB, "SETUP.md"), []byte(
        d+"go\n"+
        "type Request struct{}\n"+
        "type Response struct{}\n"+
        "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"+
        d+"\n"), 0644)

    subDir := filepath.Join(modB, "sub")
    os.MkdirAll(subDir, 0755)
    leafSetup := d+"go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+d+"\n"
    os.WriteFile(filepath.Join(subDir, "SETUP.md"), []byte(leafSetup), 0644)

    leafB := filepath.Join(subDir, "leaf-b")
    os.MkdirAll(leafB, 0755)
    os.WriteFile(filepath.Join(leafB, "SETUP.md"), []byte(leafSetup), 0644)
    os.WriteFile(filepath.Join(leafB, "ASSERT.md"), []byte(d+"go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"+d+"\n"), 0644)

    req.Args = []string{"test", leafB}
    return nil
}
```
