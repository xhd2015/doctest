# Scenario

**Feature**: wide range accepts host Go

```
min 1.0 max 99.0 -> Validate OK
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ConfigJSON = `{"go":{"min":"1.0","max":"99.0"}}`
	return nil
}
```
