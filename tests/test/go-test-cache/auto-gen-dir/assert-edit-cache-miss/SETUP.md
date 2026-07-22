# Scenario

**Feature**: a first run has completed successfully

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

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
    cfg := multiRunCfg{}
    cfg.TestDir = createTempTestProject(t, "mytest")
    cfg.ModifyFile = "simple/ASSERT.md"
    cfg.ModifyContent = doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n    if resp.Stdout != \"modified\" {\n        t.Log(\"stdout was not modified\")\n    }\n}")
    doMultiRun(t, req, cfg)
    return nil
}
```
