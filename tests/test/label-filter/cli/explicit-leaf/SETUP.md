# Scenario

**Feature**: explicit leaf paths still honor `--label`

```
doctest test <mod>/slow --label EXPR -> run or skip that leaf
```

## Steps

1. Point args at a concrete leaf directory.

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