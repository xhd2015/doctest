# Scenario

**Feature**: multi source-module gens use **one** hub suite go test (ModeHubSuite)

Production path (Phase 1): single hub plan — no multi-cmd merge / no path-shaped
executor loop required for this CLI shape.

```
outer + gotestplan/sub  →  exactly one plan:
  cd …/__hub && go test ./suite
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = createChildModuleProject(t)
	req.Args = []string{"test", "-v", "./..."}
	return nil
}
```
