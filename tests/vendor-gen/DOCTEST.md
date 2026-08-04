# Vendor-aware gen go.mod (WriteGoMod)

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Project module root** — directory with the project's `go.mod` (and optionally
  a `vendor/` tree). This is `modRoot` passed to `WriteGoMod`.
- **Vendor tree** — `modRoot/vendor/` plus `vendor/modules.txt` listing vendored
  modules, versions, optional `## explicit; go X.Y` metadata, and package paths.
- **Gen module writer (`WriteGoMod`)** — builds the nested `testcase` module's
  `go.mod` under an isolated gen directory used by doctest prepare.
- **Parent go.mod directives** — the project's `go` version line and any
  filesystem path `replace` directives that `WriteGoMod` already copies into the
  gen module (absolute-ized).
- **Vendor gomod overlay** — when a vendored module lacks `go.mod`, genDir gets
  only a synthetic `vendor-gomod-overlay/<module>/go.mod` plus
  `vendor-gomod-overlay.json` mapping
  `abs(project vendor/<mod>/go.mod) → abs(placeholder go.mod)` for `go -overlay`
  (xgo `createGoModPlaceholder` model). Package sources stay under project
  `vendor/`; project tree is never written; no package hardlink/copy into gen.

**Behaviors**

1. When `modRoot` has **no** `vendor/` directory, gen `go.mod` must **not** gain
   a mass of `replace … => …/vendor/…` lines derived from modules.txt (only the
   existing project/session/assert/parent replaces).
2. When `modRoot` has `vendor/` with usable `modules.txt`, for each module entry
   in that file gen `go.mod` includes:
   - `require <module> <version>` (version from modules.txt; zero pseudo-version
     only if version empty)
   - `replace <module> => <modRoot>/vendor/<module-path>` **always** for package
     sources (including modules that lack `go.mod`)
3. For each module whose project `vendor/<mod>/go.mod` is missing, WriteGoMod
   writes a **placeholder only** under `genDir/vendor-gomod-overlay/<mod>/go.mod`
   and records the mapping in `genDir/vendor-gomod-overlay.json`. It does **not**
   hardlink/copy package trees into gen, and does **not** point replace at a
   `vendor-bridge` shadow.
4. When `go.mod` already exists under vendor, replace still points at the project
   vendor path; **no** placeholder/overlay entry for that module.
5. Parent `go` directive and parent path replaces remain present **alongside**
   vendor requires/replaces for **other** modules (not dropped or overwritten).
6. When parent `go.mod` already has a **filesystem path** `replace` for module M
   and vendor/modules.txt also lists M, gen emits **exactly one** replace for M
   whose RHS is the parent path (absolute-ized) — **not** a second
   `replace M => …/vendor/M` (parent path replace wins; avoids dual-replace tidy
   failure). Other modules in modules.txt still get vendor require+replace.
7. When parent `go.mod` already has a **module→module** `replace` for module M
   (e.g. `replace M => M v1.3.0` version pin, or `replace M => other/mod vX`)
   and vendor/modules.txt also lists M, gen emits **exactly one** replace for M
   whose RHS is the parent module target — **not** a second
   `replace M => …/vendor/M` (parent module replace wins; same dual-replace
   class as path, seen with `github.com/gogo/protobuf`). Other modules in
   modules.txt still get vendor require+replace.
8. When modules.txt records a **private-fork** (non-local) replacement form
   `# A ver => B ver2` with packages under `vendor/A/` (not necessarily
   `vendor/B/`), gen must stay **offline-safe**: prefer
   `require A ver` + filesystem `replace A => <modRoot>/vendor/A` (or
   equivalent). Must **not** emit a bare `require B …` that forces tidy to
   network-resolve a private path without a filesystem replace for B.
   Dual-replace for the same left-hand path remains forbidden.
9. Vendor injection is **always on** when `vendor/` exists — no user flag
   (subject to the parent-replace-wins and private-fork offline-safe rules above).

## Decision Tree

```
vendor-gen/                                    [WriteGoMod + project vendor]
├── absent/                                    no vendor/ next to go.mod
│   └── no-vendor-replaces/                    no mass replace …/vendor/…
└── present/                                   vendor/ + modules.txt fixture
    ├── require-replace/                       require + replace → project vendor
    ├── placeholder-gomod/                     overlay placeholder + JSON (no package mirror)
    ├── has-gomod-no-overlay/                   existing vendor go.mod → no overlay artifacts
    ├── prefers-vendor-source/                 replace target = vendor path (+ marker)
    ├── coexists-parent/                       parent go + path replace (other mod) + vendor
    ├── parent-path-replace-wins/              parent path replace for M skips vendor/M
    ├── parent-module-replace-wins/            parent module→module replace for M skips vendor/M
    └── private-fork-offline/                  modules.txt A=>B private fork offline-safe
```

Split factor at root children: **vendor presence** (absent | present) — the
dominant product switch for this feature. Under `present/`, siblings partition
observable outcomes of vendor inject (wiring to project vendor, xgo-style
overlay placeholders for missing go.mod, no-overlay when go.mod already exists,
path preference, coexistence with parent directives, parent-replace-wins when
the same module appears in both parent replace and modules.txt — path vs
module→module as sibling leaves — and private-fork modules.txt `A => B` where
gen must not bare-require private B).

## Test Index

| Leaf | Scenario (parent design) | Expect today |
|------|--------------------------|--------------|
| `absent/no-vendor-replaces` | no mass vendor replaces | GREEN (baseline lock) |
| `present/require-replace` | require + replace → project vendor for all modules.txt entries | GREEN |
| `present/placeholder-gomod` | missing go.mod → overlay placeholder + JSON; no package mirror | GREEN |
| `present/has-gomod-no-overlay` | modules with existing go.mod get no overlay entry/JSON | GREEN |
| `present/prefers-vendor-source` | replace points at vendor; source marker visible | GREEN |
| `present/coexists-parent` | parent go + path replace (other mod) + vendor | GREEN |
| `present/parent-path-replace-wins` | parent path replace for M wins over vendor/M | GREEN |
| `present/parent-module-replace-wins` | parent module→module replace for M wins | GREEN |
| `present/private-fork-offline` | modules.txt `# A ver => B ver` private fork offline-safe | GREEN |

## How to Run

```sh
cd external/doctest-master-2026-08-04-1
doctest vet ./tests/vendor-gen
doctest test ./tests/vendor-gen --label-all
```

Library surface under test: `github.com/xhd2015/doctest/libdoc/core.WriteGoMod`.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Request drives library-level WriteGoMod against a temp project modRoot.
// Setup builds go.mod and optional vendor/modules.txt fixtures; Run calls
// WriteGoMod once and returns the gen go.mod bytes.
type Request struct {
	GenDir  string
	ModRoot string
	ModPath string
	HasMod  bool

	WithAssertReplace  bool
	AssertCacheDir     string
	WithSessionReplace bool
	SessionCacheDir    string

	// Fixture labels used by Assert (set in Setup).
	SampleModPath     string // e.g. example.com/dep
	SampleModVersion  string // e.g. v1.2.3
	NoGoModPath       string // vendored module intentionally without go.mod
	NoGoModVersion    string
	DistinctiveMarker string // unique string inside vendored package source
	ParentLocalRel    string // e.g. ./localdep when parent has path replace
	ParentLocalAbs    string // absolute localdep path after seed
	ParentGoVersion   string // e.g. 1.20 from parent go directive
	VendorRoot        string // abs path modRoot/vendor when present
}

// Response is the gen-module go.mod after WriteGoMod.
type Response struct {
	GoModContent string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	if req.GenDir == "" || req.ModRoot == "" {
		t.Fatal("Setup must set GenDir and ModRoot")
	}
	if req.ModPath == "" {
		req.ModPath = "example.com/app"
	}
	if !req.HasMod {
		req.HasMod = true
	}
	if err := core.WriteGoMod(
		req.GenDir, req.ModRoot, req.ModPath, req.HasMod,
		req.WithAssertReplace, req.AssertCacheDir,
		req.WithSessionReplace, req.SessionCacheDir,
	); err != nil {
		return &Response{}, err
	}
	data, err := os.ReadFile(filepath.Join(req.GenDir, "go.mod"))
	if err != nil {
		return &Response{}, err
	}
	return &Response{GoModContent: string(data)}, nil
}
```
