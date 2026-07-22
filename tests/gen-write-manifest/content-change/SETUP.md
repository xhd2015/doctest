# Scenario

**Feature**: desired go.mod change updates target + manifest and invalidates tidy-done

```
first WriteGoMod (parent has replace localdep)
  -> seed tidy-done
change parent go.mod (drop replace)
second WriteGoMod
  -> gen go.mod content updates
  -> manifest go.mod entry updates
  -> doctest.tidy-done removed (actual write)
```

## Preconditions

- Parent module change alters desired nested go.mod bytes.

## Steps

1. First WriteGoMod with parent replace directive.
2. Seed tidy-done; snapshot manifest entry.
3. Set `ChangeSourceGoMod` so Run rewrites parent before second WriteGoMod.

## Context

- Sibling leaves split MECE on outcome facet: content/manifest vs tidy-done.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Mode = "write-gomod-second"
	req.ModPath = "example.com/app"
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n\nreplace localdep => ./dep\n")
	firstWriteGoMod(t, req)
	seedTidyDone(t, req.GenDir)
	// Snapshot pre-change manifest line for go.mod (if present).
	man := readFileOrEmpty(manifestPath(req.GenDir))
	req.SnapManifestEntryBefore = findManifestLine(man, "go.mod")
	req.SnapManifestContentBefore = man
	req.SnapGoModContentBefore = readFileOrEmpty(filepath.Join(req.GenDir, "go.mod"))
	req.ChangeSourceGoMod = "module example.com/app\n\ngo 1.21\n"
	return nil
}
```
