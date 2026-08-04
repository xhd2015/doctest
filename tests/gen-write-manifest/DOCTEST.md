# Unified Gen Write Manifest

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Gen root** — isolated directory that holds a generated `testcase` module
  (`go.mod`, optional `go.sum`), generated leaf packages, and bookkeeping files.
  Tests always use a fresh temp gen root (never shared `mapping-gen` cache).
- **Unified write manifest** — single file at gen root, conventionally
  `doctest.gen-manifest`. Maps each path **relative to gen root** (slash-
  separated) to a **content hash** of the final on-disk bytes (post-format for
  Go sources). Includes a **version** field so format bumps force full re-hash
  misses.
- **WriteIfChanged** — shared write path used by generate writers: hash desired
  content; if `manifest[path] == hash`, skip reading the target file and skip
  rewrite; otherwise write when needed and set `manifest[path] = hash`.
- **WriteGoMod** — builds desired nested `go.mod` / copies `go.sum` into the gen
  root. Skip key is the **desired content hash via the unified manifest** only
  — not a separate policy fingerprint file.
- **Legacy gomod fingerprint** — former `doctest.gomod-fp` skip marker. After
  this feature it must not be written or required.
- **Tidy-done marker** — `doctest.tidy-done` means “`go mod tidy` already ran
  for the current module files.” It is invalidated only when `go.mod` or
  `go.sum` **actually wrote**.
- **Doctest self-module** — when `modPath == github.com/xhd2015/doctest`,
  assert/session replace flags are ineffective (no extra replace lines).
  Ineffective flag differences must not churn `go.mod` mtime or drop tidy-done
  (prior fix preserved under the unified manifest).

**Behaviors**

1. First successful module write into an empty gen root creates
   `doctest.gen-manifest` listing at least `go.mod` (and other written paths).
2. Second write with identical desired content does not rewrite `go.mod`
   (mtime stable) and does not drop `doctest.tidy-done` if present.
3. When the in-memory path→hash map is unchanged, flushing the manifest is
   content-stable (no rewrite of the manifest file itself).
4. `doctest.gomod-fp` is absent after write (never recreated as the skip key).
5. When desired `go.mod` content changes, the file and its manifest entry
   update; tidy-done is removed when `go.mod`/`go.sum` actually wrote.
6. Generic WriteIfChanged paths under the gen root (e.g. formatted leaf Go
   sources) use the same manifest: hash hit → no rewrite; hash miss → write +
   update entry.
7. For the doctest module, two WriteGoMod calls that differ only in ineffective
   assert/session replace flags keep `go.mod` mtime and tidy-done stable.

## Decision Tree

```
gen-write-manifest/
├── first-write/                         cold gen root: first WriteGoMod
│   ├── creates-manifest/                doctest.gen-manifest exists, lists go.mod
│   └── no-gomod-fp/                     doctest.gomod-fp absent
├── warm-identical/                      second write, same desired content
│   ├── go-mod-mtime-stable/             go.mod mtime unchanged
│   ├── tidy-done-retained/              doctest.tidy-done kept
│   ├── manifest-stable/                 map unchanged → manifest not rewritten
│   └── src-fp-present/                  doctest.gomod-src written (input fingerprint)
├── content-change/                      desired go.mod differs on second call
│   ├── updates-gomod-and-manifest/      go.mod + manifest entry refresh
│   ├── drops-tidy-done/                 tidy-done removed when go.mod wrote
│   └── modules-txt-invalidates/         modules.txt-only change drops tidy-done + rewrites
├── write-if-changed/                    generic relative path under gen root
│   ├── hash-hit-skips-rewrite/          same bytes → target mtime stable
│   └── hash-miss-writes/                different bytes → write + manifest entry
└── regression/                          preserve prior multi-tree fix
    └── ineffective-assert-flags/        doctest module: flag churn does not rewrite
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `first-write/creates-manifest` | First WriteGoMod creates `doctest.gen-manifest` with a `go.mod` entry |
| `first-write/no-gomod-fp` | After WriteGoMod, `doctest.gomod-fp` does not exist |
| `warm-identical/go-mod-mtime-stable` | Second identical WriteGoMod leaves `go.mod` mtime unchanged |
| `warm-identical/tidy-done-retained` | Second identical WriteGoMod keeps seeded `doctest.tidy-done` |
| `warm-identical/manifest-stable` | Second identical WriteGoMod does not rewrite the manifest file |
| `warm-identical/src-fp-present` | After warm WriteGoMod, `doctest.gomod-src` exists |
| `content-change/updates-gomod-and-manifest` | Source go.mod change updates gen go.mod content and manifest entry |
| `content-change/drops-tidy-done` | When go.mod actually writes, `doctest.tidy-done` is removed |
| `content-change/modules-txt-invalidates` | modules.txt-only change regenerates gen go.mod and drops tidy-done |
| `write-if-changed/hash-hit-skips-rewrite` | Same formatted content: target mtime stable; manifest path retained |
| `write-if-changed/hash-miss-writes` | Different content: target updates and manifest entry changes |
| `regression/ineffective-assert-flags` | Doctest module: differing ineffective assert/session flags do not churn go.mod mtime or tidy-done |

## How to Run

```sh
doctest vet ./tests/gen-write-manifest/
doctest test ./tests/gen-write-manifest/    # Classic TDD: expect RED until unified manifest lands
```

```go
import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Request drives library-level gen-root write scenarios (no CLI).
// Mode selects the measured Run operation; Setup prepares dirs / first write.
type Request struct {
	Mode string // write-gomod | write-gomod-second | write-file | write-file-second

	GenDir  string
	ModRoot string
	ModPath string
	HasMod  bool

	WithAssertReplace  bool
	AssertCacheDir     string
	WithSessionReplace bool
	SessionCacheDir    string

	// Optional second WriteGoMod flags (ineffective-assert-flags regression).
	SecondWithAssertReplace  bool
	SecondAssertCacheDir     string
	SecondWithSessionReplace bool
	SecondSessionCacheDir    string

	// If non-empty, rewrite ModRoot/go.mod to this before second WriteGoMod.
	ChangeSourceGoMod string
	// If non-empty, rewrite ModRoot/vendor/modules.txt before second WriteGoMod
	// (go.mod/go.sum may stay unchanged — must invalidate gomod-src cache).
	ChangeSourceModulesTxt string

	// Generic WriteIfChanged / WriteFormattedGo path relative to GenDir.
	RelPath         string
	FileContent     string
	SecondFileContent string // if set, second write uses this instead of FileContent

	// Per-leaf "before" snapshots for warm/second-call Asserts (Parallel-safe; not package vars).
	SnapGoModMtimeBefore      time.Time
	SnapManifestMtimeBefore   time.Time
	SnapTargetMtimeBefore     time.Time
	SnapManifestEntryBefore   string
	SnapGoModContentBefore    string
	SnapManifestContentBefore string
}

// Response captures gen-root artifacts and mtimes around the measured write.
type Response struct {
	GoModContent      string
	ManifestContent   string
	ManifestExists    bool
	GomodFpExists     bool
	TidyDoneExists    bool
	GomodSrcExists    bool // doctest.gomod-src input fingerprint
	TargetContent     string // RelPath file after Run, when applicable
	GoModMtimeBefore  time.Time
	GoModMtimeAfter   time.Time
	ManifestMtimeBefore time.Time
	ManifestMtimeAfter  time.Time
	TargetMtimeBefore time.Time
	TargetMtimeAfter  time.Time
	ManifestEntryBefore string // substring/line for RelPath or go.mod before second op
	ManifestEntryAfter  string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	resp := &Response{}
	switch req.Mode {
	case "write-gomod":
		if err := core.WriteGoMod(req.GenDir, req.ModRoot, req.ModPath, req.HasMod,
			req.WithAssertReplace, req.AssertCacheDir,
			req.WithSessionReplace, req.SessionCacheDir); err != nil {
			return resp, err
		}
	case "write-gomod-second":
		// Snapshots taken in Setup; this is the measured second WriteGoMod.
		if req.ChangeSourceGoMod != "" {
			if err := os.WriteFile(filepath.Join(req.ModRoot, "go.mod"), []byte(req.ChangeSourceGoMod), 0644); err != nil {
				return resp, err
			}
		}
		if req.ChangeSourceModulesTxt != "" {
			vendorDir := filepath.Join(req.ModRoot, "vendor")
			if err := os.MkdirAll(vendorDir, 0755); err != nil {
				return resp, err
			}
			if err := os.WriteFile(filepath.Join(vendorDir, "modules.txt"), []byte(req.ChangeSourceModulesTxt), 0644); err != nil {
				return resp, err
			}
		}
		assertRep := req.WithAssertReplace
		assertDir := req.AssertCacheDir
		sessRep := req.WithSessionReplace
		sessDir := req.SessionCacheDir
		if req.SecondAssertCacheDir != "" || req.SecondWithAssertReplace || req.SecondWithSessionReplace || req.SecondSessionCacheDir != "" {
			assertRep = req.SecondWithAssertReplace
			assertDir = req.SecondAssertCacheDir
			sessRep = req.SecondWithSessionReplace
			sessDir = req.SecondSessionCacheDir
		}
		if err := core.WriteGoMod(req.GenDir, req.ModRoot, req.ModPath, req.HasMod,
			assertRep, assertDir, sessRep, sessDir); err != nil {
			return resp, err
		}
	case "write-file":
		if err := writeGenRelFile(t, req.GenDir, req.RelPath, req.FileContent); err != nil {
			return resp, err
		}
	case "write-file-second":
		content := req.FileContent
		if req.SecondFileContent != "" {
			content = req.SecondFileContent
		}
		if err := writeGenRelFile(t, req.GenDir, req.RelPath, content); err != nil {
			return resp, err
		}
	default:
		t.Fatalf("unknown Mode: %q", req.Mode)
	}
	fillResponse(t, req, resp)
	return resp, nil
}
```
