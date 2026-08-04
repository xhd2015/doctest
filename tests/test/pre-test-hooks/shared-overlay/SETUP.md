# Scenario

**Feature**: enabled hooks accumulate in one driver-owned overlay file

```
ordered hooks -> shared generated paths -> overlay bytes decide Go flag
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
