# Scenario

**Feature**: consumer leaf that imports session and calls Once succeeds with embedded replace

```
# end-to-end consumer
doctest test leaf importing session
  -> materialize + replace
  -> Once returns valid JSON under injected DOCTEST_SESSION_ID
```

## Preconditions

- Leaf ASSERT calls `session.Once` with a small JSON object.
- Subprocess receives `DOCTEST_SESSION_ID` from doctest test runner.

## Steps

1. Create module with session-importing assert that calls Once.
2. Run `doctest test` and expect exit 0 after implementation.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t, req, "", defaultSessionAssertGo(), true)
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "-v"}
	return nil
}
```
