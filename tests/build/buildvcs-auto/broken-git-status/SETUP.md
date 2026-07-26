# Scenario

**Feature**: when git status fails, `-buildvcs=auto` errors (real fail axis)

```
# same module layout as healthy leaves, but .git/HEAD corrupted
createOriginModule -> full clone -> corrupt HEAD -> go build -buildvcs=auto
# -> error obtaining VCS status: exit status 128
```

## Preconditions

- Clone is otherwise a valid module (go.mod + main.go present).
- Only VCS metadata access is broken — models CI "dubious ownership" /
  corrupt-git class of failures, not shallow depth.

## Steps

1. Create multi-commit origin and full-clone.
2. Overwrite `.git/HEAD` with garbage so `git status --porcelain` fails.
3. `Run` builds with `-buildvcs=auto`; Assert expects the VCS status error.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	origin := createOriginModule(t)
	req.CloneDir = cloneBrokenHEAD(t, origin)
	return nil
}
```
