# Scenario

**Feature**: `doctest cache --help` documents clean and dry-run flags

```
cli.RunWithWriters -> doctest cache --help
  -> usage mentions --clean and --dry-run; exit 0
```

## Preconditions

- Cache command registered on the CLI (after implement).
- No filesystem fixture required.

## Steps

1. Set Args to `cache --help`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"cache", "--help"}
	return nil
}
```
