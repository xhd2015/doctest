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
1. Set cfg.ModifyFile = "simple/SETUP.md" to modify the leaf SETUP.md between runs
   (leaf Setup is inlined into generated source; root-only edits still regenerate,
   but leaf is the clearest cache-bust signal after inject).
2. The multi-run helper overwrites that file with new content before the second run.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    // Fully specify cfg (unified suite shares package vars across leaves).
    cfg = multiRunCfg{
        TestDir:       createTempTestProject(t, "mytest"),
        ModifyFile:    "simple/SETUP.md",
        ModifyContent: doctestGoBlock("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { req.WorkDir = \"modified-leaf-setup\"; return nil }"),
    }
    doMultiRun(t, req)
    return nil
}
```
