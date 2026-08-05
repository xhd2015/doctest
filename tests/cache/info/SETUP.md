# Scenario

**Feature**: default `doctest cache` prints cache home, root, and human sizes

```
# info path
Harness -> isolate CacheHome
  -> optional seed buckets
  -> cli.RunWithWriters(["cache"])
  -> Cache home / Doctest root / buckets / Total
```

## Preconditions

- Leaves set injectable `CacheHome` via `ensureCacheHome`.
- No `--clean` / `--dry-run` in this group.

## Steps

1. Grouping Setup ensures a temp CacheHome when not already set.
2. Leaves seed content (or leave empty) and set Args to `cache`.

## Context

- Grouping: shared CacheHome isolation default.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	return nil
}
```
