# Scenario

**Feature**: parse group for go min/max shapes

```
load test.config.json -> XgoTestConfig.Go
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.OnlyLoad = true
	return nil
}
```
