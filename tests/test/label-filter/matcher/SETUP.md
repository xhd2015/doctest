# Scenario

**Feature**: boolean label expressions match leaf label sets

```
# parse EXPR and test membership
EvalLabelExpr(expr, labels) -> match bool | parse error
```

## Steps

1. Invoke `core.EvalLabelExpr` with scenario-specific inputs in `Assert`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}

func evalLabelExpr(t *testing.T, expr string, labels []string) (bool, error) {
	t.Helper()
	return core.EvalLabelExpr(expr, labels)
}
```