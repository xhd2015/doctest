# Scenario

**Feature**: parent-internal multi-leaf gen uses unified layout markers

```
RunTest(multi-leaf parent internal, GenDir=tmp)
  -> GenDir contains suite/ and __allleaves/
  -> not classic multi-leaf-only internalCompile dump
```

## Preconditions

- Fixture from parent multi-leaf Setup; GenDir inspectable after run.

## Steps

1. Ensure cover flags off so run reaches layout dump/gen.
2. Assert unified dirs present and subject success.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	_ = d
	// Layout inspect needs a completed prepare; no coverprofile multi-pkg early exit.
	req.WithCover = false
	req.CoverPath = ""
	return nil
}
```
