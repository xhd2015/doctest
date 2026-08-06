# Scenario

**Feature**: object form `"go":{"min","max"}` loads both bounds

```
{"go":{"min":"1.18","max":"1.20"}} -> Min=1.18 Max=1.20
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ConfigJSON = `{"go":{"min":"1.18","max":"1.20"},"flags":["--unified"]}`
	return nil
}
```
