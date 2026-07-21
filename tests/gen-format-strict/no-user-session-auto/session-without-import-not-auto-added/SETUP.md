# Scenario

**Feature**: A5 — user `session` usage without import is not auto-injected by WriteFormattedGo

```
package genformatstrict
func UseSessionWithoutImport() { _ = session.Doctest{} }
  -> WriteFormattedGo
  -> no github.com/xhd2015/doctest/session import added
  -> go build fails
```

## Preconditions

- Source has no session import path.
- WantBuild true so compile failure is observed.
- Expect RED only if format path wrongly auto-adds non-stdlib/session maps; current
  `stdlibByPkgName` does not include session, so this leaf may already be GREEN as a
  regression guard that session is not added by reconcile.

## Steps

1. Set Source to package that references `session.Doctest` without importing it.
2. Assert FormattedSource lacks session import and BuildErr non-empty.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "write_format_build"
	req.WantBuild = true
	req.OutGoName = "user_session.go"
	// Intentionally missing session import; body uses session.Doctest.
	// Non-_test file so `go build .` typechecks the package.
	req.Source = `package genformatstrict

func UseSessionWithoutImport() {
	_ = session.Doctest{}
}
`
	return nil
}
```
