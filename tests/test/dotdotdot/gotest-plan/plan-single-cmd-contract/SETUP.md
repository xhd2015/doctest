# Scenario

**Feature**: `gotestmap.Plan` suite/hub modes always return **len==1**

Pure Plan API contract in Assert. Minimal CLI smoke (`--help`) so the leaf uses
inherited Run — no multi-cmd finish, no nested fixture project.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Cheap CLI so inherited Run succeeds (exit 0); contract is pure gotestmap.Plan in Assert.
	// Note: top-level "help" is unknown; use --help.
	req.Args = []string{"--help"}
	return nil
}
```
