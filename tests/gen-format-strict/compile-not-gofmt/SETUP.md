# Scenario

**Feature**: generated sources need only compile; format.Source / gofmt is not required

```
# A6
fixture with explicit imports (testing + time)
  -> generate + suite go test succeed
  -> Assert checks compile success only
  -> does NOT require gofmt-pretty output
```

## Preconditions

- Same success path as A2 for imports; different Assert contract (no gofmt).
- FixtureKind `explicit-compile`.

## Steps

1. `Op=run_fixture`, FixtureKind `explicit-compile`.
2. Assert RunErr empty; never call gofmt equality as a hard requirement.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Op = "run_fixture"
	req.FixtureKind = "explicit-compile"
	return nil
}
```
