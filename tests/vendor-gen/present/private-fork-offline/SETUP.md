# Scenario

**Bug**: modules.txt private-fork shape (`A => B`) yields bare require of private fork

```
# modules.txt (replacement form; packages under vendor by import path)
# example.com/lib v1.0.0 => example.com/fork v1.1.0
example.com/lib
... packages live under vendor/example.com/lib/ (NOT vendor/example.com/fork/)

parent go.mod:
  require example.com/lib v1.0.0
  replace example.com/lib => example.com/fork v1.1.0

  -> WriteGoMod
  -> gen go.mod offline-safe:
       require example.com/lib v1.0.0
       replace example.com/lib => <modRoot>/vendor/example.com/lib
       # NOT: bare require example.com/fork … without filesystem replace
```

## Steps

1. Override present's tiny vendor with the private-fork modules.txt shape
   (`# lib ver => fork ver`) and place packages under `vendor/example.com/lib/`.
2. Seed parent go.mod with `require example.com/lib` plus module→module
   `replace example.com/lib => example.com/fork v1.1.0` (the monorepo dual that
   surfaces network fetch of the private fork when gen requires `fork`).
3. Keep a second simple vendor module (`example.com/nogo`) so non-fork vendor
   inject still coexists.
4. Run WriteGoMod; Assert requires offline-safe graph (no bare private fork).

## Context

- Distinct from `parent-module-replace-wins`: there parent replaces **same path**
  with a version pin (`M => M vX`) and modules.txt has **no** `=>` replacement
  token. Here modules.txt records **A => B** (import path A, replacement B) and
  the failure mode is **network require of B** (private), not dual-replace tidy
  conflict alone.
- Distinct from `require-replace`: that leaf uses plain `# path version` lines
  without a non-local replacement path.
- Fixture paths stay under `example.com/…` (no real private hosts).
- Classic TDD: **RED** today — `vendorBridgeForModRoot` sets `reqPath` to the
  non-local `ReplacementPath` (`example.com/fork`) and emits
  `require example.com/fork <ver>`; with parent replacing `lib`, vendor
  filesystem replace is skipped, leaving a bare fork require that tidy must
  download.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

// Private-fork shape constants (stable fixture; not real network hosts).
const (
	forkImportPath   = "example.com/lib"  // A — import / vendor tree path
	forkImportVer    = "v1.0.0"
	forkReplacePath  = "example.com/fork" // B — private replacement module path
	forkReplaceVer   = "v1.1.0"
	forkMarker       = "VENDOR_PRIVATE_FORK_LIB_MARKER_doctest"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.ModRoot == "" || req.VendorRoot == "" {
		t.Fatal("present Setup must set ModRoot and VendorRoot")
	}

	// Labels for Assert (import path A is the offline-safe require left-hand).
	req.SampleModPath = forkImportPath
	req.SampleModVersion = forkImportVer
	req.DistinctiveMarker = forkMarker
	// Reuse ParentLocal* slots for the private replacement module B (path + ver).
	req.ParentLocalRel = forkReplacePath
	req.ParentLocalAbs = forkReplaceVer

	// --- Rewrite vendor for private-fork modules.txt shape ---
	// packages under vendor/example.com/lib/ (import path), not vendor/.../fork/
	modulesTxt := strings.Join([]string{
		"# " + forkImportPath + " " + forkImportVer + " => " + forkReplacePath + " " + forkReplaceVer,
		"## explicit; go 1.18",
		forkImportPath,
		forkImportPath + "/sub",
		"# " + noGoModPath + " " + noGoModVersion,
		"## explicit; go 1.17",
		noGoModPath,
		"",
	}, "\n")
	writeFile(t, filepath.Join(req.VendorRoot, "modules.txt"), modulesTxt)

	libDir := filepath.Join(req.VendorRoot, filepath.FromSlash(forkImportPath))
	writeFile(t, filepath.Join(libDir, "go.mod"),
		"module "+forkImportPath+"\n\ngo 1.18\n")
	writeFile(t, filepath.Join(libDir, "lib.go"),
		"package lib\n\n// "+forkMarker+"\nconst Marker = \""+forkMarker+"\"\n")
	writeFile(t, filepath.Join(libDir, "sub", "sub.go"),
		"package sub\n\nconst Y = 2\n")

	// Ensure nogo still present (present Setup already wrote it; reaffirm package).
	nogoDir := vendorModuleDir(req, noGoModPath)
	if !fileExists(filepath.Join(nogoDir, "nogo.go")) {
		writeFile(t, filepath.Join(nogoDir, "nogo.go"), "package nogo\n\nconst X = 1\n")
	}
	req.NoGoModPath = noGoModPath
	req.NoGoModVersion = noGoModVersion

	// Parent: require import path + replace to private fork (module→module).
	// This is the dual that makes bare require(fork) force a network fetch.
	body := strings.Join([]string{
		"module " + req.ModPath,
		"",
		"go 1.19",
		"",
		"require " + forkImportPath + " " + forkImportVer,
		"",
		"replace " + forkImportPath + " => " + forkReplacePath + " " + forkReplaceVer,
		"",
	}, "\n")
	seedParentGoMod(t, req, body)

	if !fileExists(filepath.Join(req.VendorRoot, "modules.txt")) {
		t.Fatal("vendor/modules.txt must exist after private-fork seed")
	}
	if !fileExists(filepath.Join(libDir, "lib.go")) {
		t.Fatal("vendor tree must place packages under import path example.com/lib")
	}
	// Must NOT place content only under fork path — offline content is at import path.
	// (Creating empty fork dir is fine; assert cares about gen go.mod graph.)
	return nil
}
```
