## Preconditions
- A first run has completed successfully.
- The auto-gen-dir parent provides multi-run Run.

## Steps
1. Set cfg.ModifyFile = "simple/ASSERT.md" to modify the leaf ASSERT.md between runs.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg.TestDir = createTempTestProject(t, "mytest")
    cfg.ModifyFile = "simple/ASSERT.md"
    cfg.ModifyContent = doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n    if resp.Stdout != \"modified\" {\n        t.Log(\"stdout was not modified\")\n    }\n}")
    doMultiRun(t, req)
    return nil
}
```
