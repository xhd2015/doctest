# Scenario

**Feature**: CLI integration for `doctest test --label`

```
doctest test <mod> --label EXPR -> filter labeled leaves -> run/skip summary
```

## Steps

1. Build standard fixture mod in leaf Setup when needed.
2. Set `req.Args` to `test`, mod path, and `--label` flags.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	return nil
}
```