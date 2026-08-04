# Scenario

**Feature**: parent go.mod content change invalidates gomod-src and rewrites gen go.mod

```
first write -> seed tidy-done
change parent go.mod only
second write
  -> gen go.mod updates (e.g. new go version / extra replace)
  -> tidy-done removed
```

## Steps

1. Fresh gen + parent go 1.21; first write; seed tidy-done.
2. Set ChangeSourceGoMod to a different parent module body (go 1.20).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = defaultModPath
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	firstWrite(t, req)
	seedTidyDone(t, req.GenDir)
	req.SnapGoModContentBefore = readFileOrEmpty(filepath.Join(req.GenDir, "go.mod"))
	// Distinct go directive so gen content must change on rebuild.
	req.ChangeSourceGoMod = "module example.com/app\n\ngo 1.20\n"
	return nil
}
```
