# Scenario

**Feature**: internal modfile is parent go.mod copy plus assert replace

```
# internal compile -modfile
WriteInternalModfile(modRoot, cacheDir) -> full parent copy + assert replace
```

## Steps

1. Create parent module with sample replace directive.
2. Call `WriteInternalModfile`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.RunKind = "internal-modfile"
	req.ModRoot = t.TempDir()
	parentGoMod := "module example.com/app\n\ngo 1.21\n\nreplace example.com/dep => ../dep\n"
	if err := os.WriteFile(filepath.Join(req.ModRoot, "go.mod"), []byte(parentGoMod), 0644); err != nil {
		t.Fatalf("write parent go.mod: %v", err)
	}
	return nil
}
```