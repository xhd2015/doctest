# Gen-root source-input go.mod cache (`doctest.gomod-src`)

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Gen root** — isolated temp directory for nested `testcase` module files and
  bookkeeping (`go.mod`, optional `go.sum`, tidy marker, fingerprints). Never
  the shared mapping-gen cache.
- **Project mod root** — parent module with `go.mod`, optional `go.sum`, and
  optional `vendor/modules.txt` whose content identity feeds the fingerprint.
- **Source fingerprint (`doctest.gomod-src`)** — multi-line input fingerprint
  (`version gomod-src=1` + hashes of project go.mod/go.sum/modules.txt +
  effective assert/session keys + modRoot/modPath/hasMod). Exact-string match
  is required for a warm hit.
- **Bridges cache (`doctest.vendor-bridges.json`)** — cached
  `[]VendorBridgeMapping` returned on warm hit without rebuild.
- **WriteGoModWithVendorBridges** — library entry that early-outs on fingerprint
  hit (gen go.mod + bridges JSON + placeholders + overlay when needed) or fully
  rebuilds on miss, then saves fingerprint + bridges.
- **Tidy-done marker** — `doctest.tidy-done`; retained on warm hit; dropped only
  when gen `go.mod`/`go.sum` actually write on a miss rebuild.
- **Vendor gomod overlay** — placeholders under `vendor-gomod-overlay/` plus
  `vendor-gomod-overlay.json` when vendored modules lack go.mod (packages stay
  in project vendor; replace still points at project vendor).

**Behaviors**

1. Cold write creates `doctest.gomod-src` and `doctest.vendor-bridges.json`
   (empty `bridges: []` when no placeholders).
2. Warm hit: no rewrite of gen go.mod / fingerprint / placeholders; tidy-done
   retained; bridges returned from cache with same `BridgeRoot` paths.
3. Source miss (parent go.mod, go.sum, modules.txt, or effective assert path)
   forces rebuild; tidy-done drops if gen mod/sum wrote.
4. Integrity miss (missing gen go.mod or a listed placeholder) forces rebuild
   and restores artifacts.

## Decision Tree

```
gomod-src-cache/                              [L2: WriteGoModWithVendorBridges]
├── first-write/                              cold gen root
│   ├── creates-src-and-bridges/              gomod-src + vendor-bridges.json exist
│   └── policy-version/                       version gomod-src=1; no legacy gomod-fp
├── warm-hit/                                 identical sources, second call
│   ├── go-mod-mtime-stable/                  gen go.mod mtime unchanged
│   ├── tidy-done-retained/                   seeded tidy-done kept
│   └── gomod-src-stable/                     fingerprint content unchanged
├── invalidate/                               source fingerprint miss
│   ├── parent-gomod/                         parent go.mod change → rebuild
│   ├── parent-gosum/                         parent go.sum create/change → miss
│   ├── modules-txt/                          modules.txt-only → gen picks module
│   └── assert-path/                          effective assert cache dir change
├── integrity/                                gen-side damage → rebuild
│   ├── missing-gen-gomod/                    delete gen go.mod → restore
│   └── missing-placeholder/                  delete placeholder → restore
└── vendor-bridges/                           bridges return value + placeholder warm
    ├── warm-returns-cached/                  same BridgeRoot on second call
    └── placeholder-mtime-stable/             placeholder not rewritten on hit
```

Split factor at root: **cache outcome class** (cold create | warm hit |
invalidate by source | integrity miss | bridges/placeholder focus). Under each
class, siblings partition the distinct observable miss/hit properties.

Overlap note: `tests/gen-write-manifest/` already has thin leaves
`warm-identical/src-fp-present` and `content-change/modules-txt-invalidates`.
This tree is the dedicated MECE matrix for gomod-src cache; those leaves remain
as regression cross-refs and are not vacuous-duplicated here beyond the needed
matrix coverage.

## Test Index

| Leaf | Scenario |
|------|----------|
| `first-write/creates-src-and-bridges` | First write creates `doctest.gomod-src` and `doctest.vendor-bridges.json` (empty bridges OK) |
| `first-write/policy-version` | Fingerprint starts with `version gomod-src=1`; `doctest.gomod-fp` absent |
| `warm-hit/go-mod-mtime-stable` | Second identical write leaves gen go.mod mtime unchanged |
| `warm-hit/tidy-done-retained` | Seeded `doctest.tidy-done` remains after warm hit |
| `warm-hit/gomod-src-stable` | Fingerprint file content equal before/after warm hit |
| `invalidate/parent-gomod` | Parent go.mod change rebuilds gen go.mod; tidy-done dropped |
| `invalidate/parent-gosum` | Parent go.sum change causes miss; tidy-done dropped when sum wrote |
| `invalidate/modules-txt` | modules.txt-only change regenerates go.mod and drops tidy-done |
| `invalidate/assert-path` | Non-doctest module: assert cache dir change forces miss rebuild |
| `integrity/missing-gen-gomod` | Delete gen go.mod then second call restores it |
| `integrity/missing-placeholder` | Delete listed placeholder then second call restores it |
| `vendor-bridges/warm-returns-cached` | Warm hit returns same BridgeRoot as cold write |
| `vendor-bridges/placeholder-mtime-stable` | With nogo module, placeholder mtime stable on warm hit |

## How to Run

```sh
doctest vet ./tests/gomod-src-cache/
doctest test ./tests/gomod-src-cache/

# regression neighbors (optional)
doctest test ./tests/gen-write-manifest/ ./tests/vendor-gen/
```

Coverage-backfill mode: product behavior already correct — **GREEN expected**.

```go
import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Request drives L2 WriteGoModWithVendorBridges scenarios on isolated temp dirs.
// Mode selects the measured Run operation; Setup prepares fixtures / first write.
type Request struct {
	Mode string // write-once | write-second

	GenDir  string
	ModRoot string
	ModPath string
	HasMod  bool

	WithAssertReplace  bool
	AssertCacheDir     string
	WithSessionReplace bool
	SessionCacheDir    string

	// Optional second-call flag overrides (assert-path invalidation).
	UseSecondFlags           bool
	SecondWithAssertReplace  bool
	SecondAssertCacheDir     string
	SecondWithSessionReplace bool
	SecondSessionCacheDir    string

	// Source mutations applied in Run before the measured second write.
	ChangeSourceGoMod      string
	ChangeSourceGoSum      string // if set, write parent go.sum to this content
	ChangeSourceModulesTxt string

	// Integrity mutations applied in Run before the measured second write.
	DeleteGenGoMod    bool
	DeletePlaceholder bool // delete SnapPlaceholderPath if set

	// Fixture: seed vendor module without go.mod (placeholder / bridges).
	SeedVendorNogo bool
	NogoModPath    string // default example.com/nogo

	// Parallel-safe per-leaf snapshots (not package vars).
	SnapGoModMtimeBefore   time.Time
	SnapGoModContentBefore string
	SnapGomodSrcBefore     string
	SnapGomodSrcMtime      time.Time
	SnapPlaceholderPath    string
	SnapPlaceholderMtime   time.Time
	SnapBridgeCount        int
	SnapBridgeRoot         string
	SnapBridgeModulePath   string
}

// Response captures bridges return value and gen-root artifacts after Run.
type Response struct {
	GoModContent         string
	GomodSrcContent      string
	GomodSrcExists       bool
	BridgesJSONExists    bool
	BridgesJSONContent   string
	GomodFpExists        bool
	TidyDoneExists       bool
	OverlayJSONExists    bool
	PlaceholderExists    bool
	PlaceholderPath      string
	BridgeCount          int
	BridgeRoots          []string
	BridgeModulePaths    []string
	GoModMtimeBefore     time.Time
	GoModMtimeAfter      time.Time
	GomodSrcMtimeBefore  time.Time
	GomodSrcMtimeAfter   time.Time
	PlaceholderMtimeBefore time.Time
	PlaceholderMtimeAfter  time.Time
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	resp := &Response{}
	switch req.Mode {
	case "write-once":
		bridges, err := core.WriteGoModWithVendorBridges(req.GenDir, req.ModRoot, req.ModPath, req.HasMod,
			req.WithAssertReplace, req.AssertCacheDir,
			req.WithSessionReplace, req.SessionCacheDir)
		fillResponse(t, req, resp, bridges)
		return resp, err
	case "write-second":
		if req.ChangeSourceGoMod != "" {
			if err := os.WriteFile(filepath.Join(req.ModRoot, "go.mod"), []byte(req.ChangeSourceGoMod), 0644); err != nil {
				return resp, err
			}
		}
		if req.ChangeSourceGoSum != "" {
			if err := os.WriteFile(filepath.Join(req.ModRoot, "go.sum"), []byte(req.ChangeSourceGoSum), 0644); err != nil {
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
		if req.DeleteGenGoMod {
			if err := os.Remove(filepath.Join(req.GenDir, "go.mod")); err != nil && !os.IsNotExist(err) {
				return resp, err
			}
		}
		if req.DeletePlaceholder && req.SnapPlaceholderPath != "" {
			if err := os.Remove(req.SnapPlaceholderPath); err != nil && !os.IsNotExist(err) {
				return resp, err
			}
		}
		assertRep := req.WithAssertReplace
		assertDir := req.AssertCacheDir
		sessRep := req.WithSessionReplace
		sessDir := req.SessionCacheDir
		if req.UseSecondFlags {
			assertRep = req.SecondWithAssertReplace
			assertDir = req.SecondAssertCacheDir
			sessRep = req.SecondWithSessionReplace
			sessDir = req.SecondSessionCacheDir
		}
		bridges, err := core.WriteGoModWithVendorBridges(req.GenDir, req.ModRoot, req.ModPath, req.HasMod,
			assertRep, assertDir, sessRep, sessDir)
		fillResponse(t, req, resp, bridges)
		return resp, err
	default:
		t.Fatalf("unknown Mode: %q", req.Mode)
	}
	return resp, nil
}
```
