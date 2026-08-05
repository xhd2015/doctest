# Scenario

**Feature**: invalid flags on `cache` fail clearly

```
cli.RunWithWriters -> doctest cache --not-a-real-flag
  -> non-zero; flag/usage error
```

## Preconditions

- No fixture tree required for flag parse errors.

## Steps

1. Leaf sets Args with an unknown flag.

## Context

- Grouping only.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	_ = req
	return nil
}
```
