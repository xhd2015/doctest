# Scenario

**Feature**: a failed first hook prevents later hooks and overlay activation

```
first hook fails -> config driver stops -> no second hook and no Go flag
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{
		{Command: []string{"first", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}},
		{Command: []string{"second", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}},
	}
	req.FailAtCall = 1
	return nil
}
```
