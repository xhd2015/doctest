# Scenario

**Feature**: parent leaf setup (dir is also a parent of child/)

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Name = "code-only"
	return nil
}
```
