# Scenario

**Feature**: runner cannot import product `internal` directly

```
go build ./suite_direct -> use of internal package example.com/app/internal/greet
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	req.Scenario = "direct"
	return nil
}
```
