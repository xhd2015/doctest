# Scenario

**Feature**: user seed plus pre_test hooks share one driver overlay (later hook wins same key)

```
user seed -> pre_test hook writes -> final Replace
```

## Preconditions

- Hooks reference `$GO_INSTRUMENT_OVERLAY_FILE` so the driver opens the shared file.
- User path is non-empty seed.

## Steps

1. Set Mode materialize (parent).
2. Leaves choose same-key vs disjoint hook writes.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreTest = []core.PreTestHook{
		{Command: []string{"hook", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"}},
	}
	req.UseMaterializeHelper = false
	return nil
}
```
