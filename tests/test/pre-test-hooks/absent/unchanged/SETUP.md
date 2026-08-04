# Scenario

**Feature**: absent `pre_test` has no generated instrumentation state

```
no pre_test -> config driver -> unchanged Go flags
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
