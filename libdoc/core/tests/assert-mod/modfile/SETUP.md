# Scenario

**Feature**: go.mod and internal modfile wiring includes assert replace when needed

```
# nested module path
WriteGoMod(..., withAssert=true) -> replace assert => cache

# internal compile path
WriteInternalModfile(parent, cache) -> parent go.mod copy + assert replace
```

## Preconditions

- Assert cache dir is materialized before modfile tests.

## Steps

1. Descendant sets `runKind` and prepares temp dirs.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	cacheDir, err := core.MaterializeAssertModule()
	if err != nil {
		t.Fatalf("materialize for modfile tests: %v", err)
	}
	req.AssertCacheDir = cacheDir
	return nil
}
```