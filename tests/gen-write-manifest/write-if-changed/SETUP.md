# Scenario

**Feature**: generic gen-relative WriteIfChanged uses the unified content-hash manifest

```
# after gen root exists (WriteGoMod)
write rel path under gen root (formatted Go source)
  -> manifest[rel] = hash(final bytes)
second write same/different content
  -> hash hit: skip rewrite (mtime stable)
  -> hash miss: write + update manifest entry
```

## Preconditions

- Gen root established with WriteGoMod so bookkeeping files live at gen root.
- Relative path is slash-separated from gen root (e.g. `leaf/leaf_test.go`).

## Steps

1. Prepare gen root via WriteGoMod.
2. Configure RelPath + FileContent; leaves perform first and/or second write.

## Context

- Production writers (`WriteFormattedGo` / leaf generate) must consult the same
  manifest. Leaves use `writeGenRelFile` which formats Go sources and expects
  the implementer to record hashes under `doctest.gen-manifest`.

```go
import (
	"testing"
)

const sampleGoA = "package leaf\n\nfunc Answer() int { return 42 }\n"
const sampleGoB = "package leaf\n\nfunc Answer() int { return 99 }\n"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ModPath = "example.com/app"
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	firstWriteGoMod(t, req)
	req.RelPath = "leaf/leaf_test.go"
	req.FileContent = sampleGoA
	return nil
}
```
