# Scenario

**Feature**: sessionmod embeds top-level session package for materialize/replace

```
# mirror assertmod
session/*.go -> embed generator -> libdoc/sessionmod/{session.go,cache_key.go}
sessionmod.ContentMD5() / RawSourceCacheKeyMD5() -> core.MaterializeSessionModule
```

## Preconditions

- Module root is four levels above this tree (`DOCTEST_ROOT/../../../..`).
- Package `github.com/xhd2015/doctest/libdoc/sessionmod` will expose
  `Content()`, `ContentMD5()`, and `RawSourceCacheKeyMD5()` after implementation.
- Classic TDD: expect RED until embed package exists.

## Steps

1. Root Setup resolves ModuleRoot.
2. Leaf sets `runKind` and Run calls sessionmod accessors.

## Context

- Consumer-facing materialize + go.mod replace is covered by
  `tests/session-inject/`, not this tree.

```go
import (
	"path/filepath"
	"testing"
)

var runKind string

func Setup(t *testing.T, req *Request) error {
	req.ModuleRoot = filepath.Join(DOCTEST_ROOT, "..", "..", "..", "..")
	return nil
}
```
