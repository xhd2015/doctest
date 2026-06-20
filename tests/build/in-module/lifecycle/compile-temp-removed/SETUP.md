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
import "testing"

func Setup(t *testing.T, req *Request) error {
	createInternalModuleProject(t)
	setupModuleEnv(t, req)
	req.Args = append(req.Args, testDir, "--gen-dir", genDir, "-v")
	return nil
}
```