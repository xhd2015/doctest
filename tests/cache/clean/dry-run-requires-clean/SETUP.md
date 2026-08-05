# Scenario

**Feature**: bare `--dry-run` without `--clean` is an error

```
doctest cache --dry-run
  -> non-zero; message: --dry-run requires --clean
```

## Preconditions

- CacheHome isolated but no seed required (flag validation only).

## Steps

1. Set Args to `cache --dry-run`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	req.Args = []string{"cache", "--dry-run"}
	return nil
}
```
