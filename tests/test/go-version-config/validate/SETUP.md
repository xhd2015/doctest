# Scenario

**Feature**: validate group runs ValidateXgoTestConfigGoVersion after load

```
load -> ValidateXgoTestConfigGoVersion(cfg, "go")
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.OnlyLoad = false
	return nil
}
```
