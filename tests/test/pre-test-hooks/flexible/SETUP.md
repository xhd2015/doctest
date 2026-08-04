# Scenario

**Feature**: pre_test argv elements expand unified placeholders as substrings

```
# mid-string tokens drive allocation + ReplaceAll substitution
pre_test command arg (contains $TOKEN) -> config driver scan/expand -> executor sees expanded argv
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
