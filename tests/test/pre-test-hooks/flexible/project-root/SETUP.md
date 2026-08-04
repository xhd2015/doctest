# Scenario

**Feature**: `$PROJECT_ROOT` expands mid-string in pre_test; unrelated `$OTHER` stays literal

```
arg "--config=$PROJECT_ROOT/cfg" + "--literal=$OTHER" -> driver expands only PROJECT_ROOT -> no overlay allocation
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{{Command: []string{
		"tool",
		"--config=$PROJECT_ROOT/cfg",
		"--literal=$OTHER",
	}}}
	return nil
}
```
