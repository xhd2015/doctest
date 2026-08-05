# Scenario

**Feature**: missing path fails with non-zero exit and Error on stderr

```
Harness -> list /nonexistent/path-for-list-doctest
  -> non-zero; stderr Error / not exist
```

## Preconditions

- Pattern points at a path that does not exist.

## Steps

1. Args = `list <absMissing>`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Under temp so we never collide with a real path; do not create it.
	missing := filepath.Join(t.TempDir(), "does-not-exist-list-root")
	req.FixtureDir = filepath.Dir(missing)
	req.Args = listArgs(nil, missing)
	return nil
}
```
