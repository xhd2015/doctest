# `doctest test` accepts `-overlay` / `--overlay` (user seed + internal merge)

## Version

0.0.2

**Layer L2 in-process** — parse via `runner.ParseTestOptions`; merge via core
overlay helpers (Replace map + single `-overlay=` GoFlags). No product-binary
e2e, no `label: e2e`. Classic TDD: feature not landed → leaves **RED** until
parse accept, abs-resolve, seed-then-later-wins, and single-flag materialize.

**Out of scope**

- `doctest build` accepting overlay (only a reject regression leaf)
- Hub tidy / `GOFLAGS` user overlay
- Full multi-mod hub e2e (L3)

# DSN (Domain Specific Notion)

## Participants

- **CLI flag parser (`ParseTestOptions`)** — accepts Go-style `-overlay FILE` and
  long form `--overlay FILE` on `doctest test` only; abs-resolves relative paths
  against process cwd (same class as profile path fields).
- **User overlay seed** — a standard Go overlay JSON (`Replace` map) supplied by
  the user; becomes the **initial** Replace map for the driver-owned overlay.
- **pre_test hooks** — ordered commands that may write the shared
  `$GO_INSTRUMENT_OVERLAY_FILE`; on the same disk-path key, **hook wins** over
  the user seed.
- **Vendor-gomod bridge normalizer** — after hooks (or without hooks), merges
  phantom `vendor/.../go.mod` → placeholder mappings; **later internal layer
  wins** the same key over seed/hook.
- **Go/xgo invocation** — receives **at most one** ordinary `-overlay=<file>`
  argument for the merged result.

## Behaviors

1. Help (covered under `tests/help/test-options`) documents `-overlay` and
   `--overlay`.
2. Both short and long flag forms parse; relative paths become absolute on
   `Options.Overlay`.
3. User overlay is **seed**; later pre_test writes and vendor-bridge merge
   overwrite the same Replace key (later wins).
4. Disjoint keys from user and internal layers all appear in the final Replace.
5. When pre_test never opens a driver overlay file: materialize one JSON
   (user seed ∪ vendor-gomod) and emit a single `-overlay=` flag.
6. No user overlay: internal pre_test / vendor behavior stays as today.
7. `doctest build` still rejects `-overlay` (test-only feature).

## Decision Tree

```text
overlay-flag/
├── parse/                         CLI accept + path + missing value
│   ├── short-form/                -overlay FILE → Options.Overlay abs
│   ├── long-form/                 --overlay FILE → Options.Overlay abs
│   ├── relative-abs/              relative path abs-resolved at parse
│   └── missing-arg/               -overlay without value → parse error
├── materialize/                   final Replace + at most one -overlay=
│   ├── user-only/                 seed only → one flag, user Replace
│   ├── user-and-hook/
│   │   ├── same-key/              hook overwrites user seed key
│   │   └── disjoint/              user key + hook key both present
│   ├── user-and-vendor/
│   │   ├── same-key/              vendor go.mod pair wins user seed key
│   │   └── disjoint/              user package key + vendor go.mod both
│   └── no-user-hook-populated/    empty user; hook-populated path unchanged
└── build-reject/                  doctest build -overlay → still error
```

Root split is **concern** (parse | materialize | build-scope). Under parse,
split by form / path / error. Under materialize, split by **which layers
contribute** (user-only | user+hook | user+vendor | no-user), then same-key vs
disjoint where both outcomes matter.

## Test Index

| Leaf | Contract | Classic TDD |
|------|----------|-------------|
| `parse/short-form` | `-overlay FILE` accepted; `Options.Overlay` set abs | RED until flag |
| `parse/long-form` | `--overlay FILE` accepted; `Options.Overlay` set abs | RED until flag |
| `parse/relative-abs` | relative path → abs (`filepath.IsAbs`, equals Abs) | RED until resolve |
| `parse/missing-arg` | `-overlay` alone → parse error naming overlay | RED until flag+arity |
| `materialize/user-only` | one `-overlay=`; Replace equals user map | RED until materialize |
| `materialize/user-and-hook/same-key` | final Replace[key] = hook value | RED until seed+hook |
| `materialize/user-and-hook/disjoint` | user key and hook key both present | RED until union |
| `materialize/user-and-vendor/same-key` | vendor mapping wins shared go.mod key | RED until vendor merge |
| `materialize/user-and-vendor/disjoint` | user package + vendor go.mod both present | RED until union |
| `materialize/no-user-hook-populated` | no user; one `-overlay=` from hook content | RED until empty-user path / GREEN once ≡ today |
| `build-reject` | `build -overlay …` still fails (unknown/rejected) | stay GREEN as reject |

## Implementer surface (L2 seams)

| API | Role |
|-----|------|
| `core.Options.Overlay string` | Abs path to user overlay JSON after parse, or empty |
| `runner.ParseTestOptions` | `String("-overlay,--overlay", &opts.Overlay)` + abs-resolve like profiles |
| `core.ApplyPreTestHooksWithUserOverlay(config, projectRoot, overlayRoot, userOverlay string, bridges, exec)` | Seed user Replace into driver `__overlay/overlay.json` **before** hooks; hooks + vendor-bridge later-win; return `PreTestHookApply` with at most one `-overlay=` |
| `core.MaterializeUserVendorOverlay(userOverlay, vendorGomodOverlay, destRoot string) (PreTestHookApply, error)` | No pre_test overlay file path: write merged user ∪ vendor JSON under destRoot, one `-overlay=` when non-empty |
| Help `testUsage` | Document `-overlay FILE` and `--overlay FILE` |
| Build path | Do **not** accept `-overlay` (keep reject) |

Empty `userOverlay` on the apply/materialize helpers must match today's
no-user behavior (pre-test-hooks tree remains the full internal suite).

### Implementer bootstrap (Classic TDD)

This tree references product symbols that do not exist yet. **First implementer
step** (stubs OK): add `Options.Overlay`, register `-overlay,--overlay` in
`parseTestOptions` (abs-resolve), and stub
`ApplyPreTestHooksWithUserOverlay` / `MaterializeUserVendorOverlay` so the
suite **compiles**. Then leaves fail on asserts (true RED) until merge/forward
behavior lands. Do not append a second `-overlay=` from
`appendOptsGoTestFlags` alone when internal already contributed one.

## How to Run

```sh
doctest vet ./tests/test/overlay-flag
doctest test ./tests/test/overlay-flag --label-all
# help leaf (separate tree):
doctest test ./tests/help/test-options
```

```go
import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cli"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/session"
)

// Mode selects the L2 surface under test.
const (
	modeParse       = "parse"
	modeMaterialize = "materialize"
	modeBuildReject = "build_reject"
)

// Request drives parse, materialize (seed ∪ internal), or build-reject.
// Classic TDD: Options.Overlay and Apply/Materialize helpers may be missing
// until the implementer lands them.
type Request struct {
	Mode string // parse | materialize | build_reject

	// --- parse ---
	ParseArgs []string // args after "test" subcommand (no "test" prefix)

	// --- materialize ---
	// UserReplace is written to a temp user overlay JSON when non-nil.
	// nil + UserOverlayEmpty means no user overlay path is passed.
	UserReplace      map[string]string
	UserOverlayEmpty bool
	// VendorReplace is written to a temp vendor-gomod-overlay.json when non-nil
	// (no-pre_test materialize path).
	VendorReplace map[string]string
	// PreTest + HookOverlays mirror tests/test/pre-test-hooks.
	PreTest      []core.PreTestHook
	HookOverlays [][]OverlayEntry
	// ActiveBridge enables xgo-style go.mod → placeholder merge after hooks.
	ActiveBridge           bool
	CreateActiveBridgeFile bool
	// When true, use MaterializeUserVendorOverlay even if PreTest is empty
	// (user-only / user+vendor no-pre_test paths).
	UseMaterializeHelper bool

	// --- build_reject ---
	BuildArgs []string
}

// OverlayEntry is one hook write into the shared overlay Replace map.
// Source is a fixture alias: "project-source", "active-vendor", "inactive-vendor".
type OverlayEntry struct {
	Source  string
	Replace string
}

type Response struct {
	// Parse
	Overlay  string // Options.Overlay after successful parse
	ParseErr string

	// Materialize
	OverlayFile    string
	GoFlags        []string
	OverlayReplace map[string]string
	OverlayFlagN   int // count of args with prefix -overlay=
	// Fixture paths for asserts (filled by Run).
	ProjectRoot        string
	ProjectSource      string
	ActiveVendorSource string
	ActiveVendorGoMod  string
	ActiveBridgeSource string

	// Shared
	ExitCode int
	ErrMsg   string
	Stdout   string
	Stderr   string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	switch req.Mode {
	case modeParse:
		opts, _, err := runner.ParseTestOptions(req.ParseArgs)
		if err != nil {
			resp.ExitCode = 1
			resp.ParseErr = err.Error()
			resp.ErrMsg = err.Error()
			return resp, nil
		}
		// Options.Overlay is the product field the implementer adds.
		resp.Overlay = opts.Overlay
		return resp, nil
	case modeBuildReject:
		var buf bytes.Buffer
		err := cli.RunWithWriter(&buf, req.BuildArgs)
		resp.Stdout = buf.String()
		if err != nil {
			resp.ExitCode = 1
			resp.ErrMsg = err.Error()
			resp.Stderr = err.Error() + "\n"
			return resp, nil
		}
		return resp, nil
	case modeMaterialize:
		return runMaterialize(t, req, resp)
	default:
		resp.ExitCode = 1
		resp.ErrMsg = "unknown Mode: " + req.Mode
		return resp, nil
	}
}

func runMaterialize(t *testing.T, req *Request, resp *Response) (*Response, error) {
	t.Helper()
	fixtureRoot := t.TempDir()
	projectRoot := filepath.Join(fixtureRoot, "project")
	generatedRoot := filepath.Join(fixtureRoot, "generated")
	if err := os.MkdirAll(projectRoot, 0o755); err != nil {
		return nil, err
	}
	activeModule := "example.com/active"
	inactiveModule := "example.com/inactive"
	activeVendorRoot := filepath.Join(projectRoot, "vendor", filepath.FromSlash(activeModule))
	activePlaceholderGoMod := filepath.Join(generatedRoot, "vendor-gomod-overlay", filepath.FromSlash(activeModule), "go.mod")
	activeVendorSource := filepath.Join(activeVendorRoot, "pkg", "active.go")
	activeVendorGoMod := filepath.Join(activeVendorRoot, "go.mod")
	projectSource := filepath.Join(projectRoot, "pkg", "project.go")
	inactiveVendorSource := filepath.Join(projectRoot, "vendor", filepath.FromSlash(inactiveModule), "pkg", "inactive.go")
	for _, source := range []string{activeVendorSource, projectSource, inactiveVendorSource} {
		if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(source, []byte("package fixture\n"), 0o644); err != nil {
			return nil, err
		}
	}
	if req.CreateActiveBridgeFile {
		if err := os.MkdirAll(filepath.Dir(activePlaceholderGoMod), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(activePlaceholderGoMod, []byte("module example.com/active\n\ngo 1.19\n"), 0o644); err != nil {
			return nil, err
		}
	}
	var bridges []core.VendorBridgeMapping
	if req.ActiveBridge {
		bridges = append(bridges, core.VendorBridgeMapping{
			ModulePath:         activeModule,
			OriginalVendorRoot: activeVendorRoot,
			BridgeRoot:         activePlaceholderGoMod,
		})
	}

	resp.ProjectRoot = projectRoot
	resp.ProjectSource = projectSource
	resp.ActiveVendorSource = activeVendorSource
	resp.ActiveVendorGoMod = activeVendorGoMod
	resp.ActiveBridgeSource = activePlaceholderGoMod

	// Expand fixture aliases in Replace maps so leaves can write SETUP before
	// TempDir paths exist (project-source, active-vendor, inactive-vendor,
	// active-vendor-gomod). Value sentinel $ACTIVE_BRIDGE → placeholder go.mod.
	aliasPaths := map[string]string{
		"project-source":      projectSource,
		"active-vendor":       activeVendorSource,
		"inactive-vendor":     inactiveVendorSource,
		"active-vendor-gomod": activeVendorGoMod,
	}
	valueSentinels := map[string]string{
		"$ACTIVE_BRIDGE": activePlaceholderGoMod,
	}
	userReplace := expandReplaceMap(req.UserReplace, aliasPaths, valueSentinels)
	vendorReplace := expandReplaceMap(req.VendorReplace, aliasPaths, valueSentinels)

	userPath := ""
	if !req.UserOverlayEmpty && userReplace != nil {
		p, err := writeOverlayJSON(filepath.Join(fixtureRoot, "user-overlay.json"), userReplace)
		if err != nil {
			return nil, err
		}
		userPath = p
	}

	vendorPath := ""
	if vendorReplace != nil {
		p, err := writeOverlayJSON(filepath.Join(generatedRoot, core.VendorGomodOverlayJSON), vendorReplace)
		if err != nil {
			return nil, err
		}
		vendorPath = p
	}

	cfg := &core.XgoTestConfig{PreTest: req.PreTest}
	call := 0
	exec := func(workDir string, command []string) error {
		call++
		if call <= len(req.HookOverlays) {
			for i := 0; i+1 < len(command); i++ {
				if command[i] != "--overlay-file" {
					continue
				}
				overlay := struct {
					Replace map[string]string `json:"Replace"`
				}{Replace: map[string]string{}}
				if data, err := os.ReadFile(command[i+1]); err == nil && len(data) > 0 {
					if err := json.Unmarshal(data, &overlay); err != nil {
						return err
					}
					if overlay.Replace == nil {
						overlay.Replace = map[string]string{}
					}
				}
				alias := map[string]string{
					"active-vendor":   activeVendorSource,
					"project-source":  projectSource,
					"inactive-vendor": inactiveVendorSource,
				}
				for _, entry := range req.HookOverlays[call-1] {
					source := alias[entry.Source]
					if source == "" {
						return os.ErrInvalid
					}
					overlay.Replace[source] = entry.Replace
				}
				data, err := json.Marshal(overlay)
				if err != nil {
					return err
				}
				if err := os.WriteFile(command[i+1], data, 0o644); err != nil {
					return err
				}
			}
		}
		return nil
	}

	var apply core.PreTestHookApply
	var err error
	if req.UseMaterializeHelper || (len(req.PreTest) == 0 && (userPath != "" || vendorPath != "")) {
		// No pre_test driver file: materialize user ∪ vendor into one JSON/flag.
		apply, err = core.MaterializeUserVendorOverlay(userPath, vendorPath, generatedRoot)
	} else {
		// pre_test path: seed user Replace before hooks; vendor bridges after.
		apply, err = core.ApplyPreTestHooksWithUserOverlay(cfg, projectRoot, generatedRoot, userPath, bridges, exec)
	}
	if err != nil {
		resp.ExitCode = 1
		resp.ErrMsg = err.Error()
		// Still surface partial apply when present.
	}
	resp.OverlayFile = apply.OverlayFile
	resp.GoFlags = append([]string(nil), apply.GoFlags...)
	resp.OverlayFlagN = countOverlayFlags(apply.GoFlags)
	if apply.OverlayFile != "" {
		data, readErr := os.ReadFile(apply.OverlayFile)
		if readErr != nil {
			if resp.ErrMsg == "" {
				return nil, readErr
			}
		} else if len(data) > 0 {
			var overlay struct {
				Replace map[string]string `json:"Replace"`
			}
			if uerr := json.Unmarshal(data, &overlay); uerr != nil {
				if resp.ErrMsg == "" {
					return nil, uerr
				}
			} else {
				resp.OverlayReplace = overlay.Replace
			}
		}
	}
	return resp, nil
}

// expandReplaceMap rewrites known alias keys to absolute fixture paths and
// known value sentinels (e.g. $ACTIVE_BRIDGE). Unknown keys/values pass through.
func expandReplaceMap(in map[string]string, keyAlias, valueSentinel map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		key := k
		if abs, ok := keyAlias[k]; ok {
			key = abs
		}
		val := v
		if rep, ok := valueSentinel[v]; ok {
			val = rep
		}
		out[key] = val
	}
	return out
}

func writeOverlayJSON(path string, replace map[string]string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	body, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{Replace: replace})
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return "", err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path, nil
	}
	return abs, nil
}

func countOverlayFlags(flags []string) int {
	n := 0
	for _, f := range flags {
		if strings.HasPrefix(f, "-overlay=") || f == "-overlay" || strings.HasPrefix(f, "--overlay") {
			n++
		}
	}
	return n
}
```
