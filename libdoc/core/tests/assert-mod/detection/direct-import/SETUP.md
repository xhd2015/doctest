# Scenario

**Feature**: direct assert import in ASSERT block is detected

```
import "github.com/xhd2015/doctest/assert" -> CasesImportAssertPackage true
```

## Steps

1. Build case with direct assert import path.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, req *Request) error {
	req.Cases = []core.TreeCase{makeCaseWithAssertImport("")}
	return nil
}
```