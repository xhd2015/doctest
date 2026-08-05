# Scenario

**Feature**: empty match soft-exits with `no tests` on stderr and empty stdout

```
Harness -> empty dir (no DOCTEST.md)
  -> list empty/...
  -> exit 0; stderr "no tests"; stdout empty
```

## Preconditions

- Temp directory with no DOCTEST.md roots.

## Steps

1. Create empty temp dir.
2. Args = `list <dir>/...`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	dir := t.TempDir()
	req.FixtureDir = dir
	req.Args = listArgs(nil, filepath.ToSlash(dir)+"/...")
	return nil
}
```
