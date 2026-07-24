# Scenario

**Feature**: invalid `--label` expression fails before running tests

```
doctest test --label 'slow &&' -> parse error, non-zero exit
```

## Steps

1. Run test against fixture mod with invalid expression.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	_ = t
	return nil
}
```