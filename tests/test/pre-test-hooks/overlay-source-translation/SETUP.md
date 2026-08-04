# Scenario

**Bug**: caller-owned overlay source keys must follow only bridges active in this generated run

```
pre_test hooks -> shared original-path overlay -> active bridge metadata -> generated Go overlay
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	return nil
}
```
