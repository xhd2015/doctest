# Scenario

**Feature**: hooks request only the artifacts whose whole-arg overlay placeholders appear

```
# exact whole-arg tokens (flexible mid-string is under flexible/)
pre_test command -> exact placeholder scan -> generated artifact selection -> hook executor
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
