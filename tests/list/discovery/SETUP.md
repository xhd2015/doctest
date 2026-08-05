# Scenario

**Feature**: discover doctest roots by path / `...` patterns and emit body lines

```
# fixtures under t.TempDir
Harness -> write DOCTEST roots + leaves
  -> cli.RunWithWriters(["list", absPatterns...])
  -> body lines (sorted) + summary | soft no tests | errors
```

## Preconditions

- Each leaf builds its own fixture under `t.TempDir()`.
- Patterns are absolute so process cwd is not required.

## Steps

1. Grouping Setup is a no-op.
2. Leaves create fixtures and set Args.

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
