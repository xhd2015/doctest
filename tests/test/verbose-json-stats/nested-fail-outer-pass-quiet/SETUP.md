# Scenario

**Feature**: quiet path already uses json counts; nested FAIL ( must not affect PASS (2/2)

```
# same outer tree as -v case, without -v
doctest test --no-color <outer>
  -> go test -json accounting
  -> final summary PASS (2/2), exit 0
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Same outer fixture shape as `nested-fail-outer-pass-v` (2 outer passes + nested fail child).

## Steps

1. Create outer nested-fail-outer-pass fixture.
2. Isolate caches for the harness invocation.
3. Run `doctest test --no-color <outer>` (no `-v`).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	outer := createNestedFailOuterPassTree(t, req.Bin)
	req.Args = []string{"test", "--no-color", outer}
	req.Env = append(req.Env, isolateRunEnv(t)...)
	return nil
}
```
