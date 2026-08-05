# Scenario

**Feature**: bare `...` pattern is rejected (same message family as test/vet)

```
Harness -> list ...
  -> non-zero; stderr bare '...' not supported
```

## Preconditions

- No fixture required.

## Steps

1. Args = `list ...`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"list", "..."}
	return nil
}
```
