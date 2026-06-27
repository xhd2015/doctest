# Scenario

**Feature**: cases without assert import path return false

```
import "fmt" only -> CasesImportAssertPackage false
```

## Steps

1. Build case without assert import.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, req *Request) error {
	req.Cases = []core.TreeCase{makeCaseWithoutAssertImport()}
	return nil
}
```