# Scenario

**Feature**: list help and top-level dispatch surfaces

```
# list usage
cli.RunWithWriters -> doctest list --help -> usage on stdout

# top-level includes list
cli.RunWithWriters -> doctest --help -> command list includes list

# unknown flag
cli.RunWithWriters -> doctest list --not-a-real-flag -> non-zero + stderr
```

## Preconditions

- No fixture trees required for help/dispatch leaves.
- In-process only.

## Steps

1. Grouping Setup is a no-op.
2. Leaves set `req.Args` for the help or error variant.

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
