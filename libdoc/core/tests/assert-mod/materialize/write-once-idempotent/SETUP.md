# Scenario

**Feature**: second MaterializeAssertModule call does not rewrite cache files

```
# warm cache
call 1 -> call 2 -> assert.go MD5 unchanged
```

## Steps

1. Call `MaterializeAssertModule` twice.
2. Snapshot MD5 before second call in `req` via package vars.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

var beforeDigest [16]byte

func Setup(t *testing.T, req *Request) error {
	dir1, err := core.MaterializeAssertModule()
	if err != nil {
		t.Fatalf("first materialize: %v", err)
	}
	beforeDigest = snapshotMD5(t, filepath.Join(dir1, "assert.go"))
	_, err = core.MaterializeAssertModule()
	if err != nil {
		t.Fatalf("second materialize: %v", err)
	}
	return nil
}
```