# Parent-leaf unified package (one dir, one package)

## Version
0.0.1

A directory that is both a leaf (ASSERT) and an intermediate (has a child
leaf) must generate `setup.go` and `leaf.go` in the **same** Go package.

```
code-only/           parent leaf + intermediate
└── child/           child leaf
```

`code-only` sanitizes to `code_only`. Both generated files in that dir
must be `package code_only`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	Name string
}

type Response struct {
	Name string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	return &Response{Name: req.Name}, nil
}
```
