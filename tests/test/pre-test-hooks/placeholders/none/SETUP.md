# Scenario

**Feature**: a hook without overlay placeholders still runs

```
plain hook command -> config driver -> executor receives unchanged command
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{{Command: []string{"tool", "prepare"}}}
	return nil
}
```
