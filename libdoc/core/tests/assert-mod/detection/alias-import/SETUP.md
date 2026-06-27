# Scenario

**Feature**: aliased assert import path is detected

```
import outputassert "github.com/xhd2015/doctest/assert" -> true
```

## Steps

1. Build case with aliased assert import.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, req *Request) error {
	req.Cases = []core.TreeCase{makeCaseWithAssertImport("outputassert")}
	return nil
}
```