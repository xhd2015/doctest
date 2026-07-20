# Scenario

**Feature**: second WriteGoMod with identical desired content is a warm skip

```
# warm path
first WriteGoMod -> seed tidy-done -> force old mtimes
second WriteGoMod (same sources/flags)
  -> go.mod mtime unchanged
  -> tidy-done retained
  -> manifest map unchanged → manifest file not rewritten
  -> still no doctest.gomod-fp
```

## Preconditions

- First WriteGoMod already populated gen root.
- Desired go.mod content is unchanged between calls.

## Steps

1. Prepare fresh dirs and run first WriteGoMod in Setup.
2. Seed `doctest.tidy-done` and force old mtimes on `go.mod` / manifest.
3. Run Mode `write-gomod-second` as the measured call.

## Context

- Sibling leaves assert MECE stability properties of the same warm outcome.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Mode = "write-gomod-second"
	req.ModPath = "example.com/app"
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	firstWriteGoMod(t, req)
	seedTidyDone(t, req.GenDir)
	snapshotGoModMtime(t, req)
	snapshotManifestMtime(t, req)
	return nil
}
```
