# Scenario

**Feature**: internal import with --gen-dir inside module dumps review copy without nested go.mod

```
# --gen-dir inside module: temp compile + dump copy
doctest test <tests> --gen-dir <module>/_gen -> .doctest_run_* compile -> dump to _gen
```

## Preconditions

- Temp module with `internal/greet` imported in harness `Run()`.
- `--gen-dir` is `<moduleRoot>/_gen` (inside parent module).

## Steps

1. Create internal-import temp module.
2. Run `doctest test <tests> --gen-dir <moduleRoot>/_gen -v`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	createInternalModuleProject(t)
	setupModuleEnv(t, req)
	req.Args = append(req.Args, testDir, "--gen-dir", genDir, "-v")
	return nil
}
```