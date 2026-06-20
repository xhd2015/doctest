# Scenario

**Feature**: internal import without --gen-dir compiles in temp dir under module root

```
# no --gen-dir: internal import triggers temp compile only
doctest test <tests> -v -> .doctest_run_* under moduleRoot -> go test -> temp removed
```

## Preconditions

- Temp module with `internal/greet` imported in harness `Run()`.
- No `--gen-dir` flag (cache/temp compile path only).

## Steps

1. Create internal-import temp module.
2. Run `doctest test <tests> -v` without `--gen-dir`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	createInternalModuleProject(t)
	_ = genDir
	setupModuleEnv(t, req)
	req.Args = append(req.Args, testDir, "-v")
	return nil
}
```