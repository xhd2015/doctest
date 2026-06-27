# Assert Module Core Helpers — Direct Tests

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Import scanner** — `CasesImportAssertPackage` inspects parsed import paths in
  SETUP/ASSERT Go blocks for `github.com/xhd2015/doctest/assert` (including aliases).
- **Materializer** — `MaterializeAssertModule` writes content-addressed cache at
  `$CACHE/doctest/assert-mod/<md5>/` (write-once).
- **Go mod writer** — `WriteGoMod` appends `replace assert => <cache>` for nested modules.
- **Modfile builder** — copies parent `go.mod` and appends assert replace for internal compile.

**Behaviors**

- Detection true only when assert import path matches exactly (alias name ignored).
- Materialize creates `assert.go` + `go.mod` with `go 1.18` on first call; skips rewrite.
- Nested go.mod and internal modfile include assert replace only when cases import assert.

## Decision Tree

```
assert-mod/
├── detection/                         [CasesImportAssertPackage]
│   ├── direct-import/               assert import in ASSERT → true
│   ├── alias-import/                aliased path → true
│   └── absent-import/               no assert path → false
├── materialize/                     [MaterializeAssertModule cache]
│   ├── creates-cache-layout/        first call writes assert.go + go.mod
│   └── write-once-idempotent/       second call does not rewrite files
└── modfile/                         [go.mod / -modfile wiring]
    ├── nested-gomod-assert-replace/ WriteGoMod adds assert replace
    └── internal-modfile-assert-replace/ parent copy + assert replace
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `detection/direct-import` | Assert import in ASSERT Go block detected |
| `detection/alias-import` | Aliased assert import path detected |
| `detection/absent-import` | No assert import returns false |
| `materialize/creates-cache-layout` | Materialize creates cache dir with correct layout |
| `materialize/write-once-idempotent` | Second materialize leaves files unchanged |
| `modfile/nested-gomod-assert-replace` | WriteGoMod includes assert replace when requested |
| `modfile/internal-modfile-assert-replace` | Internal modfile copies parent go.mod + assert replace |

## How to Run

```sh
doctest vet ./libdoc/core/tests/assert-mod/
doctest test ./libdoc/core/tests/assert-mod/    # expect RED before implementation
```

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

type Request struct {
	Cases		[]core.TreeCase
	GenDir		string
	ModRoot		string
	ModPath		string
	CacheDir	string
	ModfilePath	string
	WithAssertReplace	bool
	AssertCacheDir	string
}
type Response struct {
	Detected	bool
	CacheDir	string
	AssertGoMD5Before	[16]byte
	AssertGoMD5After	[16]byte
	GoModContent	string
	ModfilePath	string
	ModfileContent	string
}
func Run(t *testing.T, req *Request) (*Response, error) {
	switch runKind {
	case "detect":
		return &Response{Detected: core.CasesImportAssertPackage(req.Cases, req.ModPath)}, nil
	case "materialize":
		cacheDir, err := core.MaterializeAssertModule()
		if err != nil {
			return nil, err
		}
		return &Response{CacheDir: cacheDir}, nil
	case "write-gomod":
		if err := core.WriteGoMod(req.GenDir, req.ModRoot, req.ModPath, true, req.WithAssertReplace, req.AssertCacheDir); err != nil {
			return nil, err
		}
		data, err := os.ReadFile(filepath.Join(req.GenDir, "go.mod"))
		if err != nil {
			return nil, err
		}
		return &Response{GoModContent: string(data)}, nil
	case "internal-modfile":
		path, err := core.WriteInternalModfile(req.ModRoot, req.AssertCacheDir)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return &Response{ModfilePath: path, ModfileContent: string(data)}, nil
	default:
		t.Fatalf("unknown runKind: %s", runKind)
		return nil, nil
	}
}
```