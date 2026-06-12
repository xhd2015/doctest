## Expected
- DOCTEST.md, SETUP.md, basic/SETUP.md, and basic/ASSERT.md all exist.
- The root SETUP.md contains Request and Response types and a Run function.

```go
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    dir := ""
    for _, env := range req.Env {
        if strings.HasPrefix(env, "TEST_TREE_DIR=") {
            dir = env[len("TEST_TREE_DIR="):]
            break
        }
    }
    if dir == "" {
        t.Fatal("TEST_TREE_DIR not set by SETUP")
    }

    for _, name := range []string{
        "DOCTEST.md",
        "SETUP.md",
        filepath.Join("basic", "SETUP.md"),
        filepath.Join("basic", "ASSERT.md"),
    } {
        p := filepath.Join(dir, name)
        if _, err := os.Stat(p); err != nil {
            t.Fatalf("missing file %s: %v", p, err)
        }
    }

    rootSetup, err := os.ReadFile(filepath.Join(dir, "SETUP.md"))
    if err != nil {
        t.Fatal(err)
    }
    content := string(rootSetup)
    if !strings.Contains(content, "type Request struct") {
        t.Fatal("root SETUP.md missing Request type")
    }
    if !strings.Contains(content, "func Run") {
        t.Fatal("root SETUP.md missing Run function")
    }
}
```
