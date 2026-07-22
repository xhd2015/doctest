# Scenario

**Feature**: WriteGoMod appends assert replace for nested testcase module

```
# legacy nested module with assert cases
WriteGoMod(genDir, modRoot, modPath, hasMod=true, withAssert=true)
```

## Steps

1. Create temp gen dir and parent module root.
2. Call WriteGoMod with assert replace enabled.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
    req.RunKind = "write-gomod"
	req.GenDir = t.TempDir()
	req.ModRoot = t.TempDir()
	if err := os.WriteFile(filepath.Join(req.ModRoot, "go.mod"), []byte("module example.com/app\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write parent go.mod: %v", err)
	}
	req.WithAssertReplace = true
	return nil
}
```