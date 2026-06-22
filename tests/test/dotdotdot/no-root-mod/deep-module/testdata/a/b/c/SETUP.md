# Scenario

**Feature**: tests for c

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Name = "c"
	return nil
}
```
