# Scenario

**Feature**: without vendor/, WriteGoMod does not inject vendor replaces

```
# project has go.mod only — no vendor/
modRoot (go.mod)
  -> WriteGoMod
  -> gen go.mod: project replace + optional parent path replaces
  -> NO mass replace … => …/vendor/…
```

## Preconditions

- `modRoot` has a valid `go.mod` and **no** `vendor/` directory.
- Parent may still carry ordinary path replaces; those are not vendor injection.

## Steps

1. Create temp modRoot + genDir.
2. Seed parent go.mod with module path, go directive, and one local path replace
   (proves non-vendor replaces still allowed).
3. Do not create `vendor/`.

## Context

- Scenario 7 baseline: absence of vendor means no modules.txt-driven replaces.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	prepareFreshDirs(t, req)
	req.ParentLocalRel = "./localdep"
	localAbs := filepath.Join(req.ModRoot, "localdep")
	req.ParentLocalAbs = localAbs
	if err := os.MkdirAll(localAbs, 0755); err != nil {
		t.Fatalf("mkdir localdep: %v", err)
	}
	writeFile(t, filepath.Join(localAbs, "go.mod"), "module example.com/localdep\n\ngo 1.19\n")
	body := "module " + req.ModPath + "\n\ngo 1.19\n\nreplace example.com/localdep => ./localdep\n"
	seedParentGoMod(t, req, body)
	// Ensure vendor does not exist.
	if fileExists(filepath.Join(req.ModRoot, "vendor")) {
		t.Fatal("absent branch must not create vendor/")
	}
	return nil
}
```
