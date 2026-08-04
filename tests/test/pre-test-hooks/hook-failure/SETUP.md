# Scenario

**Feature**: hook errors stop the build before any Go command is assembled

```
hook exits non-zero -> config driver -> actionable error, no test build
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
