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
1. Set cfg.ModifyFile = "SETUP.md" to modify the root SETUP.md between runs.
2. The Run function will overwrite SETUP.md with new content before the second run.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg.TestDir = createTempTestProject(t, "mytest")
    cfg.ModifyFile = "SETUP.md"
    cfg.ModifyContent = doctestGoBlock("import \"testing\"\n\ntype Request struct{ Args []string; WorkDir string }\ntype Response struct{ ExitCode int; Stdout string; Stderr string }\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\nfunc Run(t *testing.T, req *Request) (*Response, error) { return &Response{Stdout: \"modified\"}, nil }")
    doMultiRun(t, req)
    return nil
}
```
