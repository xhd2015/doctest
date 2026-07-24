# Scenario

**Feature**: nested go.mod contains replace for embedded assert cache

```
# assert import in ASSERT.md
doctest test --gen-dir outside -> generated go.mod -> replace assert => <cache>
```

## Preconditions

- Leaf ASSERT imports `github.com/xhd2015/doctest/assert`.
- Outside gen-dir triggers nested testcase module generation.

## Steps

1. Create public module with assert-importing leaf.
2. Run `doctest test <tests> --gen-dir <outsideGenDir> -v`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createPublicModuleProject(t, req, "", defaultAssertAssertGo())
	req.OutsideGenDir = filepath.Join(t.TempDir(), "generated")
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "--gen-dir", req.OutsideGenDir, "-v"}
	return nil
}
```