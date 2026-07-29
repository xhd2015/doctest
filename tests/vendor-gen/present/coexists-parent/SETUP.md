# Scenario

**Feature**: parent go directive + parent path replaces coexist with vendor replaces

```
parent go.mod:
  module example.com/app
  go 1.20
  replace example.com/localdep => ./localdep

+ vendor/modules.txt modules
  -> gen go.mod has ALL of:
       go 1.20
       replace example.com/app => <modRoot>
       replace example.com/localdep => <abs localdep>
       require+replace example.com/dep (vendor)
```

## Steps

1. Rebuild parent go.mod with go 1.20 and local path replace (override present default).
2. Keep tiny vendor fixture from present Setup.
3. Assert all three families of directives after WriteGoMod.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Override parent go.mod after present Setup's default seed.
	req.ParentLocalRel = "./localdep"
	localAbs := filepath.Join(req.ModRoot, "localdep")
	req.ParentLocalAbs = localAbs
	if err := os.MkdirAll(localAbs, 0755); err != nil {
		t.Fatalf("mkdir localdep: %v", err)
	}
	writeFile(t, filepath.Join(localAbs, "go.mod"), "module example.com/localdep\n\ngo 1.20\n")
	writeFile(t, filepath.Join(localAbs, "local.go"), "package localdep\n")
	body := "module " + req.ModPath + "\n\ngo 1.20\n\nreplace example.com/localdep => ./localdep\n"
	seedParentGoMod(t, req, body)
	if req.ParentGoVersion != "1.20" {
		t.Fatalf("expected ParentGoVersion 1.20, got %q", req.ParentGoVersion)
	}
	// Vendor fixture must still exist from present Setup.
	if !fileExists(filepath.Join(req.VendorRoot, "modules.txt")) {
		t.Fatal("vendor/modules.txt must remain after parent go.mod rewrite")
	}
	return nil
}
```
