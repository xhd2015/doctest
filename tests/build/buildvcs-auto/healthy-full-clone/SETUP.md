# Scenario

**Feature**: full clone + `-buildvcs=auto` succeeds and stamps VCS metadata

```
# full history clone of synthetic origin
createOriginModule -> git clone file://origin full -> go build -buildvcs=auto
```

## Preconditions

- Origin module has multiple commits (full clone has history depth > 1).

## Steps

1. Create multi-commit origin.
2. Full-clone into a temp work tree.
3. Run sets `CloneDir` to that clone; `Run` builds with `-buildvcs=auto`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	origin := createOriginModule(t)
	req.CloneDir = cloneFull(t, origin)
	return nil
}
```
