# Scenario

**Feature**: host above max fails validation

```
max 1.0 -> error contains "> 1.0"
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ConfigJSON = `{"go":{"max":"1.0"}}`
	return nil
}
```
