# Scenario

**Feature**: compile temp is removed after test while gen-dir dump persists

```
# lifecycle: .doctest_run_* gone after test; --gen-dir dump remains
doctest test --gen-dir <module>/_gen -> temp compile -> remove temp -> keep dump
```

## Preconditions

- Temp module with `internal/greet` imported in harness `Run()`.
- `--gen-dir` is `<moduleRoot>/_gen` (inside parent module).

## Steps

1. Create internal-import temp module.
2. Run `doctest test <tests> --gen-dir <moduleRoot>/_gen -v`.
3. Verify no `.doctest_run_*` under `moduleRoot` and dump still exists.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createInternalModuleProject(t, d, req)
	setupModuleEnv(t, req)
	req.Args = append(req.Args, req.TestDir, "--gen-dir", req.GenDir, "-v")
	return nil
}
```