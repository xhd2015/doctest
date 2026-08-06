# Scenario

**Feature**: importer outside `http/` cannot import `http/internal/leaf` directly

```
go build ./suite_direct -> use of internal package … not allowed
```

## Steps

1. Set scenario `direct`.

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
