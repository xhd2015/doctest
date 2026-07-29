# Scenario

**Feature**: multi source-module gens use hub suite go test (ModeHubSuite)

```
outer + gotestplan/sub  →  cd …/__hub && go test ./suite
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
