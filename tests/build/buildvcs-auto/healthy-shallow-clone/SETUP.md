# Scenario

**Feature**: shallow clone (`--depth 1`) + `-buildvcs=auto` also succeeds

```
# shallow clone is NOT the fail axis for -buildvcs=auto
createOriginModule -> git clone --depth 1 file://origin shallow -> go build -buildvcs=auto
```

## Preconditions

- Origin has multiple commits; clone is verified shallow via
  `git rev-parse --is-shallow-repository`.

## Steps

1. Create multi-commit origin.
2. Shallow-clone (`--depth 1`) into a temp work tree.
3. `Run` builds with `-buildvcs=auto` (same as full-clone leaf).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	origin := createOriginModule(t)
	req.CloneDir = cloneShallow(t, origin)
	return nil
}
```
