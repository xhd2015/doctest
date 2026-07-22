# Scenario

**Feature**: internal-compile with --gen-dir dumps review copy without nested go.mod

```
# internal + assert, --gen-dir inside module
doctest test --gen-dir _gen -> temp compile -> dump copy without go.mod
```

## Preconditions

- Fixture imports internal greet and assert.
- `--gen-dir` is `<moduleRoot>/_gen`.

## Steps

1. Copy internal+assert fixture.
2. Run `doctest test <tests> --gen-dir <genDir> -v`.

```go
import (
"github.com/xhd2015/doctest/session"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createInternalAssertProject(t, d, req)
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "--gen-dir", filepath.Join(req.ModuleRoot, "_gen"), "-v"}
	return nil
}
```