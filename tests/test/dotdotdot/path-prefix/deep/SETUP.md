# Scenario

**Feature**: testdata/ has go.mod at root with a DOCTest tree at group/subgroup/

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- testdata/ has go.mod at root with a DOCTest tree at group/subgroup/.

## Steps
1. Set WorkDir to the testdata/ directory.
2. Run `doctest test -v ./group/subgroup/...`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.WorkDir = "./testdata"
    req.Args = []string{"test", "-v", "./group/subgroup/..."}
    return nil
}
```
