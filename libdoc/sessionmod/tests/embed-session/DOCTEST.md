# Sessionmod Embed Package — Build-Time Tests

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Top-level session package** (`github.com/xhd2015/doctest/session`) — production
  sources that get embedded for offline/self-contained doctest trees.
- **sessionmod** (`libdoc/sessionmod`) — embeds session package sources and exposes
  content bytes, MD5, and a raw-source cache key (mirror of `libdoc/assertmod`).
- **Embed cache** — materialization target
  `$UserCacheDir/doctest/session-mod/<md5>/` with `go.mod` module
  `github.com/xhd2015/doctest/session` plus source files (owned by core materialize;
  this tree checks embed accessors).

**Behaviors**

- `Content()` / embedded file bytes match on-disk generated session source blob.
- `ContentMD5()` equals MD5 of the embedded bytes on disk.
- `RawSourceCacheKeyMD5()` is a stable non-empty hex key used for cache directory
  naming (same role as assertmod).

## Decision Tree

```
embed-session/
└── embed-package/                 [libdoc/sessionmod accessors]
    ├── content-md5-matches/       M1: ContentMD5 matches embedded file hash
    └── cache-key-nonempty/        M2: RawSourceCacheKeyMD5 is stable hex
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `embed-package/content-md5-matches` | M1 — `ContentMD5()` matches MD5 of embedded session source file |
| `embed-package/cache-key-nonempty` | M2 — `RawSourceCacheKeyMD5()` is non-empty lowercase hex |

## How to Run

```sh
doctest vet ./libdoc/sessionmod/tests/embed-session/
doctest test ./libdoc/sessionmod/tests/embed-session/   # expect RED until sessionmod lands
go test ./libdoc/sessionmod/
```

```go
import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/sessionmod"
)

type Request struct {
	RunKind	string
	ModuleRoot string
}

type Response struct {
	ContentMD5      string
	FileMD5         string
	PackageCacheKey string
	ContentLen      int
	Err             error
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.RunKind {
	case "sessionmod-md5":
		// Prefer generated session.go (or equivalent) next to embed.go.
		candidates := []string{
			filepath.Join(req.ModuleRoot, "libdoc", "sessionmod", "session.go"),
			filepath.Join(req.ModuleRoot, "libdoc", "sessionmod", "once.go"),
		}
		var data []byte
		var err error
		for _, p := range candidates {
			data, err = os.ReadFile(p)
			if err == nil {
				break
			}
		}
		if err != nil {
			return &Response{Err: fmt.Errorf("read embedded session source: %w", err)}, err
		}
		sum := md5.Sum(data)
		return &Response{
			FileMD5:    fmt.Sprintf("%x", sum),
			ContentMD5: sessionmod.ContentMD5(),
			ContentLen: len(sessionmod.Content()),
		}, nil
	case "sessionmod-cache-key":
		return &Response{
			PackageCacheKey: sessionmod.RawSourceCacheKeyMD5(),
			ContentLen:      len(sessionmod.Content()),
		}, nil
	default:
		t.Fatalf("unknown req.RunKind: %s", req.RunKind)
		return nil, nil
	}
}
```
