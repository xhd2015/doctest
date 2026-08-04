# Scenario

**Feature**: no pre-test configuration preserves the existing test invocation

```
project config -> doctest config driver -> no hook execution or overlay artifact
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
