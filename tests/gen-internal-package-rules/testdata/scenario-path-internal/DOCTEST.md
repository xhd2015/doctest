# Kind A subject — scenario path contains `internal`

Minimal tree: leaf path `http/internal/post-succeeds` with public-only Run.
Used by outer RED leaf `kind-a-scenario-path-fails`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct{}

type Response struct {
	OK bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	_ = req
	// No product internal import — failure must come from gen packaging only.
	return &Response{OK: true}, nil
}
```
