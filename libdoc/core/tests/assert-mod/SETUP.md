# Scenario

**Feature**: core package helpers detect assert imports and materialize cached module

```
# CasesImportAssertPackage scans parsed imports
SETUP/ASSERT Go blocks -> true when path == github.com/xhd2015/doctest/assert

# MaterializeAssertModule writes cache
embedded bytes -> $CACHE/doctest/assert-mod/<md5>/
```

## Preconditions

- Package `github.com/xhd2015/doctest/libdoc/core` is importable.
- `runKind` is set by descendant Setup to select helper under test.

## Steps

1. Descendant configures `req` fields and `runKind` for its scenario.

## Context

- `makeCaseWithAssertImport` / `makeCaseWithoutAssertImport` build minimal `TreeCase` values.
- `snapshotMD5` reads file and returns MD5 digest for idempotency checks.

```go
import (
	"crypto/md5"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/assertmod"
	"github.com/xhd2015/doctest/libdoc/core"
)

const assertImportPath = "github.com/xhd2015/doctest/assert"

func Setup(t *testing.T, req *Request) error {
	req.ModPath = "example.com/app"
	return nil
}
func makeCaseWithAssertImport(alias string) core.TreeCase {
	path := assertImportPath
	name := ""
	if alias != "" {
		name = alias
	}
	return core.TreeCase{
		Name:	"leaf",
		Path:	"leaf",
		AssertFile: core.AssertDocument{
			GoBlock: core.GoBlock{
				Imports: []core.ImportSpec{{Name: name, Path: path}},
			},
		},
	}
}
func makeCaseWithoutAssertImport() core.TreeCase {
	return core.TreeCase{
		Name:	"leaf",
		Path:	"leaf",
		AssertFile: core.AssertDocument{
			GoBlock: core.GoBlock{
				Imports: []core.ImportSpec{{Path: "fmt"}},
			},
		},
	}
}
func snapshotMD5(t *testing.T, path string) [16]byte {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var digest [16]byte
	copy(digest[:], h.Sum(nil))
	return digest
}
func readFileString(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
func assertCacheLayoutCore(t *testing.T, cacheDir string) {
	t.Helper()
	assertGo := filepath.Join(cacheDir, "assert.go")
	goMod := filepath.Join(cacheDir, "go.mod")
	if _, err := os.Stat(assertGo); err != nil {
		t.Fatalf("missing cached assert.go: %v", err)
	}
	modData := readFileString(t, goMod)
	if !strings.Contains(modData, "module "+assertImportPath) || !strings.Contains(modData, "go 1.18") {
		t.Fatalf("unexpected cached go.mod:\n%s", modData)
	}
	src := readFileString(t, assertGo)
	if !strings.Contains(src, "package assert") {
		t.Fatalf("cached assert.go missing package assert")
	}
	legacyNames, err := assertmod.LegacyV1Filenames()
	if err != nil {
		t.Fatalf("legacy_v1 filenames: %v", err)
	}
	for _, name := range legacyNames {
		legacyPath := filepath.Join(cacheDir, "legacy_v1", name)
		if _, err := os.Stat(legacyPath); err != nil {
			t.Fatalf("missing cached legacy_v1/%s: %v", name, err)
		}
	}
	legacyV2Names, err := assertmod.LegacyV2Filenames()
	if err != nil {
		t.Fatalf("legacy_v2 filenames: %v", err)
	}
	for _, name := range legacyV2Names {
		legacyPath := filepath.Join(cacheDir, "legacy_v2", name)
		if _, err := os.Stat(legacyPath); err != nil {
			t.Fatalf("missing cached legacy_v2/%s: %v", name, err)
		}
	}
}
```