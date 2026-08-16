# Parent-leaf unified package (one dir, one package)

## Version
0.0.1

## DSN (Domain Specific Notion)

### Participants

- **Unified generator** — writes `setup.go` (intermediate) and `leaf.go` (parent leaf) under the same dir.
- **Parent leaf** — `code-only/` has ASSERT and a child leaf.
- **Child leaf** — `code-only/child/`.

### Behaviors

- Both generated files in `code-only/` are `package code_only`.
- `doctest test` runs both leaves; no `found packages testcase … and code_only`.

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
