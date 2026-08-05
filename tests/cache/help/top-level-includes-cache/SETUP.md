# Scenario

**Feature**: top-level `doctest --help` lists the `cache` command

```
cli.RunWithWriters -> doctest --help
  -> Commands section includes cache
```

## Preconditions

- Top-level usage string is the product contract for discoverability.

## Steps

1. Set Args to `--help`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--help"}
	return nil
}
```
