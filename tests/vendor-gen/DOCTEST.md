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
- **Placeholder go.mod** — a minimal `module <path>` + `go <ver>` file written
  under a vendored module directory when that module lacks a real `go.mod`, so
  the replace target is a valid Go module (xgo-style).

**Behaviors**

1. When `modRoot` has **no** `vendor/` directory, gen `go.mod` must **not** gain
   a mass of `replace … => …/vendor/…` lines derived from modules.txt (only the
   existing project/session/assert/parent replaces).
2. When `modRoot` has `vendor/` with usable `modules.txt`, for each module entry
   in that file gen `go.mod` includes:
   - `require <module> <version>` (version from modules.txt; zero pseudo-version
     only if version empty)
   - `replace <module> => <modRoot>/vendor/<module-path>` (or replacement path
     rules aligned with xgo when modules.txt records a non-local `=>` replace)
3. For each such replace target, if `go.mod` is missing under the vendored path,
   WriteGoMod ensures a **placeholder** `go.mod` exists there (`module` + `go`).
4. Resolution prefers vendored sources: gen `go.mod` replace targets point at the
   project vendor tree (observable on the written file; optional read of
   distinctive fixture content at the replace target).
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
    ├── require-replace/                       require + replace per modules.txt
    ├── placeholder-gomod/                     placeholder go.mod when missing
    ├── prefers-vendor-source/                 replace target = vendor path (+ marker)
    ├── coexists-parent/                       parent go + path replace (other mod) + vendor
    ├── parent-path-replace-wins/              parent path replace for M skips vendor/M
    ├── parent-module-replace-wins/            parent module→module replace for M skips vendor/M
    └── private-fork-offline/                  modules.txt A=>B private fork offline-safe
```

Split factor at root children: **vendor presence** (absent | present) — the
dominant product switch for this feature. Under `present/`, siblings partition
observable outcomes of the vendor bridge (wiring, placeholders, path preference,
coexistence with existing parent directives, parent-replace-wins when the
same module appears in both parent replace and modules.txt — path vs
module→module as sibling leaves — and private-fork modules.txt `A => B` where
gen must not bare-require private B).

## Test Index

| Leaf | Scenario (parent design) | Expect today |
|------|--------------------------|--------------|
| `absent/no-vendor-replaces` | no mass vendor replaces | GREEN (baseline lock) |
| `present/require-replace` | require + replace from modules.txt | GREEN (vendor inject landed) |
| `present/placeholder-gomod` | placeholder go.mod for module without one | GREEN |
| `present/prefers-vendor-source` | replace points at vendor; source marker visible | GREEN |
| `present/coexists-parent` | parent go + path replace (other mod) + vendor | GREEN |
| `present/parent-path-replace-wins` | 1.1 — parent path replace for M wins over vendor/M | GREEN (path skip landed) |
| `present/parent-module-replace-wins` | 1.1b — parent module→module replace for M wins | GREEN (module skip landed) |
| `present/private-fork-offline` | modules.txt `# A ver => B ver` private fork offline-safe | **RED** until require stays on A + FS vendor replace |

## How to Run

```sh
cd external/doctest-master-2026-07-28-1
# product binary from this tree
go build -o /tmp/doctest-pfork ./cmd/doctest
/tmp/doctest-pfork vet ./tests/vendor-gen/
/tmp/doctest-pfork test ./tests/vendor-gen/ --label-all
# Classic TDD: present/private-fork-offline RED until implement offline-safe private-fork graph
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
