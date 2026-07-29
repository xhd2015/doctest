# Scenario

**Feature**: gen go.mod has no mass of replace … => …/vendor/… when vendor/ missing

```
WriteGoMod(no vendor/)
  -> gen go.mod without /vendor/ replace targets from modules.txt
  -> still has replace project => modRoot and parent path replace
```

## Steps

1. Inherit absent Setup (parent go.mod + localdep path replace, no vendor).
2. Run WriteGoMod; assert no vendor-path replaces.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.ModRoot == "" || req.GenDir == "" {
		t.Fatal("parent must prepare ModRoot and GenDir")
	}
	if fileExists(filepath.Join(req.ModRoot, "vendor")) {
		t.Fatal("vendor/ must be absent for this leaf")
	}
	return nil
}
```
