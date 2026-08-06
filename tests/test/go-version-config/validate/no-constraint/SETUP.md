# Scenario

**Feature**: empty config has no go constraints

```
{} -> Validate OK
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.ConfigJSON = `{"flags":["--unified"]}`
	return nil
}
```
