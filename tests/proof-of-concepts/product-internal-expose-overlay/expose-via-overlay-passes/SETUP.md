# Scenario

**Feature**: runner imports virtual product expose package via overlay

```
go run -overlay=… ./suite_expose -> hello-from-app-internal
expose path not on disk under app/
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Scenario = "expose-overlay"
	return nil
}
```
