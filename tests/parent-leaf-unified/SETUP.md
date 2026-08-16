# Scenario

**Feature**: shared root setup for parent-leaf package clash fixture

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Name == "" {
		req.Name = "root"
	}
	return nil
}
```
