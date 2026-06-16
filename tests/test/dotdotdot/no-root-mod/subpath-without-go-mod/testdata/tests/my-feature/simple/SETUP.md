# Scenario

**Feature**: simple leaf

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    t.Logf("setup: %s", req.Name)
    return nil
}
```
