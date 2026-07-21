# Scenario

**Bug**: nested FAIL (0/1) in the -v stream deflates outer suite Passed

```
# outer: pass_leaf + nested_fail_ok (nested intentional fail, outer Assert passes)
doctest test -v --no-color <outer>
  -> go test -v dumps nested stdout containing "FAIL (0/1)"
  -> stats must ignore nested FAIL ( text
  -> final summary PASS (2/2), exit 0
```

## Preconditions

- Root Setup has set `req.Bin` and a generous outer timeout.
- Outer fixture has exactly two leaves, both Assert-pass; nested child fails.

## Steps

1. Create outer nested-fail-outer-pass fixture (uses `req.Bin` for nested shell-out).
2. Isolate GOCACHE / leaf-cache / cache-home for the harness invocation.
3. Run `doctest test -v --no-color <outer>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	outer := createNestedFailOuterPassTree(t, req.Bin)
	req.Args = []string{"test", "-v", "--no-color", outer}
	req.Env = append(req.Env, isolateRunEnv(t)...)
	return nil
}
```
