# Scenario

**Feature**: first MaterializeAssertModule call creates cache directory layout

```
# cold cache dir for content hash
MaterializeAssertModule -> assert.go + go.mod with go 1.18
```

## Steps

1. Remove expected cache dir if present.
2. Call `MaterializeAssertModule` once.

```go
import (
	"os"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	cacheDir, err := core.MaterializeAssertModule()
	if err == nil {
		os.RemoveAll(cacheDir)
	}
	return nil
}
```