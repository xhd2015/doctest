# Scenario

**Feature**: file placeholder pre-creates an empty overlay JSON file

```
file placeholder -> config driver -> zero-byte overlay file -> no Go flag
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.PreTest = []core.PreTestHook{{Command: []string{"tool", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}}}
	return nil
}
```
