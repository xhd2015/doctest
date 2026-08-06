# Scenario

**Feature**: string form `"go":"1.18"` sets min only

```
{"go":"1.18"} -> Min=1.18 Max=""
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ConfigJSON = `{"go":"1.18"}`
	return nil
}
```
