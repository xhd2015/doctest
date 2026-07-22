# Scenario

**Feature**: SIGINT during internal-import temp compile removes `.doctest_run_*`

```
# internal import -> .doctest_run_* under moduleRoot
doctest test -v <tests> -> writeCases -> SIGINT (Ctrl-C) -> compile temp removed

# mechanism
signal handler calls ctx.Close() before exit so in-module compile temp is not left behind
```

## Preconditions

- Temp module with `internal/greet` imported in harness `Run()`.
- In-place fixture `./testdata/` (40 leaves) colocated with this test, copied to a temp module root.
- `InterruptDuringWriteCases` sends SIGINT after `leaf15_test.go` appears in stderr.

## Steps

1. Copy `./testdata/` into a temp module root.
2. Run `doctest test <tests> -v` and send SIGINT during `writeCases`.
3. Verify no `.doctest_run_*` remains under `moduleRoot`.

```go
import (
"testing"
"github.com/xhd2015/doctest/session"
)
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createInternalModuleProjectWithLeaves(t, d, req)
	setupModuleEnv(t, req)
	req.InterruptDuringWriteCases = true
	req.InterruptTriggerLeaf = 15
	req.Args = append(req.Args, req.TestDir, "-v")
	return nil
}
```