# Scenario

**Feature**: aliased assert import survives assembly into generated test source

```
# import outputassert "github.com/xhd2015/doctest/assert"
doctest test --gen-dir outside -> generated leaf_test.go keeps alias
```

## Preconditions

- Leaf ASSERT uses `outputassert` alias for assert package path.

## Steps

1. Create public module with aliased assert import in ASSERT.md.
2. Run `doctest test <tests> --gen-dir <outsideGenDir> -v`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t, req, "", aliasedAssertAssertGo())
	req.OutsideGenDir = filepath.Join(t.TempDir(), "generated")
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "--gen-dir", req.OutsideGenDir, "-v"}
	return nil
}
```