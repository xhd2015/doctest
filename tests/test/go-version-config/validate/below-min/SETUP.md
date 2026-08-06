# Scenario

**Feature**: host below min fails validation

```
min 99.0 -> error contains "< 99.0"
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ConfigJSON = `{"go":{"min":"99.0"}}`
	return nil
}
```
