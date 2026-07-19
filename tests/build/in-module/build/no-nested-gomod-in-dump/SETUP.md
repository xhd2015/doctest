# Scenario

**Feature**: doctest build dumps test files without nested go.mod when internal imports detected

```
# build + internal import + --gen-dir inside module
doctest build <tests> --gen-dir <module>/_gen -> dump leaf_test.go, no go.mod
```

## Preconditions

- Temp module with `internal/greet` imported in harness `Run()`.
- `--gen-dir` is `<moduleRoot>/_gen` (inside parent module).

## Steps

1. Create internal-import temp module.
2. Run `doctest build <tests> --gen-dir <moduleRoot>/_gen -v`.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createInternalModuleProject(t, d)
	setupModuleEnv(t, req)
	req.Args = append(req.Args, testDir, "--gen-dir", genDir, "-v")
	return nil
}
```