# Assert Embed Generator — Build-Time Tests

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Embed script** (`script/embed-assert`) — reads `assert/*.go`, skips `*_test.go`,
  concatenates into single `assert.go` with stable filename ordering.
- **Generated assert.go** — committed input for `//go:embed` in `libdoc/assertmod`.
- **Embed package** (`libdoc/assertmod`) — exposes `Content()` bytes and `ContentMD5()` hash.

**Behaviors**

- Script produces one concatenated Go source file with `package assert`.
- Repeated runs yield identical bytes (deterministic sort).
- `go generate` refreshes `assert.go`; `ContentMD5()` matches on-disk file hash.

## Decision Tree

```
embed-assert/
├── generator/                      [embed script output shape]
│   ├── single-concatenated-file/   A1: one file, package assert, no test sources
│   └── deterministic-bytes/        A2: two runs produce identical output
└── embed-package/                  [libdoc/assertmod accessors]
    └── md5-matches-assert-go/      A3: ContentMD5 matches assert.go hash
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `generator/single-concatenated-file` | A1 — single `assert.go`; no `*_test.go` content |
| `generator/deterministic-bytes` | A2 — two script runs produce identical bytes |
| `embed-package/md5-matches-assert-go` | A3 — `ContentMD5()` matches embedded file MD5 |

## How to Run

```sh
doctest vet ./libdoc/assertmod/tests/embed-assert/
doctest test ./libdoc/assertmod/tests/embed-assert/    # expect RED before implementation
```

```go
import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/assertmod"
)

type Request struct {
	ModuleRoot	string
	AssertDir	string
	OutputPath	string
	SecondRun	bool
}
type Response struct {
	OutputBytes	[]byte
	OutputMD5	string
	SecondRunMD5	string
	ContentMD5	string
	FileMD5		string
	Err		error
}
func Run(t *testing.T, req *Request) (*Response, error) {
	switch runKind {
	case "embed-script":
		return runEmbedScript(t, req)
	case "assertmod-md5":
		assertGo := filepath.Join(req.ModuleRoot, "libdoc", "assertmod", "assert.go")
		data, err := os.ReadFile(assertGo)
		if err != nil {
			return &Response{Err: err}, err
		}
		fileSum := md5.Sum(data)
		return &Response{
			FileMD5:	fmt.Sprintf("%x", fileSum),
			ContentMD5:	assertmod.ContentMD5(),
		}, nil
	default:
		t.Fatalf("unknown runKind: %s", runKind)
		return nil, nil
	}
}
```