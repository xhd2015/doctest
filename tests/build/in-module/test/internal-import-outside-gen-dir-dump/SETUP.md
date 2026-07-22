# Scenario

**Feature**: internal import with --gen-dir outside module dumps review copy from temp compile

```
# --gen-dir outside module: temp compile in-module, dump copied outside
doctest test <tests> --gen-dir <outside> -> .doctest_run_* under moduleRoot -> dump outside
```

## Preconditions

- Temp module with `internal/greet` imported in harness `Run()`.
- `--gen-dir` is outside the parent module (separate temp directory).

## Steps

1. Create internal-import temp module.
2. Run `doctest test <tests> --gen-dir <outsideGenDir> -v`.

```go
import (
"github.com/xhd2015/doctest/session"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createInternalModuleProject(t, d, req)
	req.OutsideGenDir = filepath.Join(t.TempDir(), "outside_dump")
	setupModuleEnv(t, req)
	req.Args = append(req.Args, req.TestDir, "--gen-dir", req.OutsideGenDir, "-v")
	return nil
}
```