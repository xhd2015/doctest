# Generic pre-test hooks and shared Go overlay

## Version

0.0.2

# DSN (Domain Specific Notion)

## Participants

- **Project test configuration** — optionally declares an ordered `pre_test`
  list. Each hook has a command argument list; any argv element may contain
  unified placeholders as substrings.
- **Doctest config driver** — reads the generic hook list, allocates only the
  requested generated artifacts (when any arg *contains* an overlay token),
  substitutes placeholders mid-string, and invokes hooks before assembling the
  Go test command. Production wiring uses a durable mapping-gen root as
  `overlayRoot` / generated root — not suite-only temp paths.
- **Hook executor** — receives each already-expanded command and the project
  root as its working directory. It may write replacement mappings to the
  shared overlay file.
- **Shared overlay** — one generated directory and one zero-byte JSON file per
  config evaluation when requested. Hooks accumulate mappings in that file.
- **Go/xgo invocation** — receives exactly one ordinary `-overlay=<file>`
  argument only when the shared file became non-empty.

## Behaviors

1. Without `pre_test`, existing config application is unchanged and creates no
   instrumentation artifacts.
2. Unified placeholder vocabulary (same replace rules for pre_test command args
   and, in production, xgo `args`): `$PROJECT_ROOT`,
   `$GO_INSTRUMENT_OVERLAY_DIR`, `$GO_INSTRUMENT_OVERLAY_FILE`. Expansion is
   **flexible substring** replacement (`Contains` / `ReplaceAll`) inside any
   argv element — not whole-arg-only. Values are config interpolation, not a
   process-environment contract; unrelated `$OTHER` tokens are left untouched.
3. Need overlay dir/file if **any** pre_test command arg **contains** the
   respective token. A referenced directory is allocated once beneath the
   generated workspace; a referenced file is allocated once, pre-created empty,
   and is usable even when no hook references the directory alone.
4. Hooks run in declaration order from the project root and share the same
   substituted paths.
5. A failed hook stops the build. After successful hooks, the file's byte size
   controls whether the driver contributes a standard `-overlay` Go argument.

## Decision Tree

```text
pre-test-hooks/
├── absent/
│   └── unchanged/                 no pre_test → no artifacts or Go flag
├── placeholders/                  exact whole-arg overlay tokens (superset-valid)
│   ├── none/                      hook runs, no instrumentation allocation
│   ├── directory-only/            shared directory only (exact arg)
│   └── file-only-empty/           pre-created file, no overlay Go flag (exact)
├── flexible/                      mid-string / unified expand vocabulary
│   ├── mid-string-file-empty/     --overlay=$GO_INSTRUMENT_OVERLAY_FILE
│   ├── mid-string-directory/      --dir=$GO_INSTRUMENT_OVERLAY_DIR/...
│   └── project-root/              $PROJECT_ROOT mid-string; $OTHER untouched
├── shared-overlay/
│   ├── populated/                 hook writes file → one -overlay flag
│   └── two-hooks-in-order/        same paths, declaration-order execution
├── overlay-source-translation/    xgo-style phantom go.mod merge after hooks
│   ├── active-bridge/              package key stays; go.mod→placeholder merged
│   ├── source-and-value-unchanged/ project + package keys/values stay; go.mod pair added
│   ├── stale-bridge-inactive/      on-disk placeholder alone cannot merge
│   ├── inactive-or-missing/        package keys stay put with active metadata
│   ├── no-active-mappings/         no bridges → vendor package keys unchanged
│   └── after-all-hooks/            hooks merge first; then go.mod pair merge
└── hook-failure/
    └── stops-before-build/        surfaced failure, no Go flag
```

The root split is whether the configuration has hooks. When hooks exist, the
next factors are which placeholders appear and whether they are exact whole-arg
tokens or mid-string interpolations. Remaining leaves distinguish overlay
activation, ordering, and failure.

## Test Index

| Leaf | Contract |
|---|---|
| `absent/unchanged` | No `pre_test` keeps the existing invocation unchanged. |
| `placeholders/none` | A normal hook runs without allocating overlay artifacts. |
| `placeholders/directory-only` | Exact directory placeholder receives one generated directory only. |
| `placeholders/file-only-empty` | Exact file placeholder receives an empty JSON file and contributes no Go flag. |
| `flexible/mid-string-file-empty` | Mid-string `$GO_INSTRUMENT_OVERLAY_FILE` allocates file, substitutes into the arg, no Go flag while empty. |
| `flexible/mid-string-directory` | Mid-string `$GO_INSTRUMENT_OVERLAY_DIR` allocates dir and substitutes inside the arg (suffix preserved). |
| `flexible/project-root` | Mid-string `$PROJECT_ROOT` expands to the absolute project root; unrelated `$OTHER` is not expanded. |
| `shared-overlay/populated` | A populated file contributes one absolute `-overlay` argument. |
| `shared-overlay/two-hooks-in-order` | Hooks share paths and execute in declaration order. |
| `overlay-source-translation/active-bridge` | Active bridges keep package overlay keys on project vendor and merge phantom go.mod→placeholder. |
| `overlay-source-translation/source-and-value-unchanged` | Project/package source keys and replacement values stay; only go.mod pair is merged. |
| `overlay-source-translation/stale-bridge-inactive` | An on-disk placeholder without this run's metadata has no effect. |
| `overlay-source-translation/inactive-or-missing` | Package overlay keys for active/inactive vendors stay put (no package rewrite). |
| `overlay-source-translation/no-active-mappings` | No bridge metadata leaves vendor package keys unchanged. |
| `overlay-source-translation/after-all-hooks` | Ordered hooks merge first; then go.mod placeholder pair is appended. |
| `hook-failure/stops-before-build` | Non-zero hook failure is returned before test build invocation. |

## Implementer surface (generic L2 seam)

The leaves describe an in-process contract for the core helper. The helper
is intentionally generic: it has no knowledge of SPL, KMS, DB, or overlay JSON
content.

```go
type PreTestHook struct {
	Command []string `json:"command"`
}

type PreTestHookExecutor func(workDir string, command []string) error

type PreTestHookApply struct {
	OverlayDir  string
	OverlayFile string
	GoFlags     []string
}

func ApplyPreTestHooks(
	config *XgoTestConfig,
	projectRoot string,
	generatedRoot string,
	exec PreTestHookExecutor,
) (PreTestHookApply, error)
```

Generated-workspace callers use a bridge-aware variant. Mappings are produced
while the current run builds vendor inject (xgo-style): packages stay under
project `vendor/`; `BridgeRoot` is the absolute placeholder **go.mod** path under
`vendor-gomod-overlay/`. Scanning a retained hardlink `vendor-bridge` tree is not
valid.

```go
type VendorBridgeMapping struct {
	ModulePath         string
	OriginalVendorRoot string // project vendor/<mod>
	BridgeRoot         string // absolute placeholder go.mod path
}

func ApplyPreTestHooksWithVendorBridges(
	config *XgoTestConfig,
	projectRoot string,
	overlayRoot string,
	bridges []VendorBridgeMapping,
	exec PreTestHookExecutor,
) (PreTestHookApply, error)
```

`XgoTestConfig` carries `PreTest []PreTestHook`. The surrounding build runner
calls this helper before package arguments are finalized, appending `GoFlags`
before packages. This tree exercises the core contract using an injected
executor; it does not need a real child process. Flexible expansion of the
unified placeholder set is part of `ApplyPreTestHooks` (and the same rules
apply to xgo `args` expansion in production — pure `args` edges may stay L1).

## How to Run

```sh
doctest vet ./tests/test/pre-test-hooks
doctest test ./tests/test/pre-test-hooks --label-all
```

```go
import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

type Request struct {
	PreTest                []core.PreTestHook
	WriteOverlay           bool
	FailAtCall             int
	ActiveBridge           bool
	CreateActiveBridgeFile bool // write placeholder go.mod under vendor-gomod-overlay
	HookOverlays           [][]OverlayEntry
}

type OverlayEntry struct {
	Source  string
	Replace string
}

type Response struct {
	OverlayDirExists     bool
	OverlayFileExists    bool
	OverlayFileSize      int64
	OverlayDir           string
	OverlayFile          string
	GoFlags              []string
	Calls                [][]string
	WorkDirs             []string
	ErrMsg               string
	OverlayReplace       map[string]string
	ProjectRoot          string
	ActiveVendorSource   string // project vendor/.../pkg/active.go (package overlay key)
	ActiveVendorGoMod    string // project vendor/.../go.mod (phantom overlay source)
	ActiveBridgeSource   string // absolute placeholder go.mod (BridgeRoot / overlay value)
	ProjectSource        string
	InactiveVendorSource string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
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
	// xgo-style: BridgeRoot is placeholder go.mod under vendor-gomod-overlay (not a package tree).
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

	resp := &Response{
		ProjectRoot:          projectRoot,
		ActiveVendorSource:   activeVendorSource,
		ActiveVendorGoMod:    activeVendorGoMod,
		ActiveBridgeSource:   activePlaceholderGoMod,
		ProjectSource:        projectSource,
		InactiveVendorSource: inactiveVendorSource,
	}
	cfg := &core.XgoTestConfig{PreTest: req.PreTest}
	call := 0
	exec := func(workDir string, command []string) error {
		call++
		resp.WorkDirs = append(resp.WorkDirs, workDir)
		resp.Calls = append(resp.Calls, append([]string(nil), command...))
		if req.FailAtCall == call {
			return errors.New("hook failed")
		}
		if req.WriteOverlay || call <= len(req.HookOverlays) {
			for i := 0; i+1 < len(command); i++ {
				if command[i] == "--overlay-file" {
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
					if call <= len(req.HookOverlays) {
						for _, entry := range req.HookOverlays[call-1] {
							source := map[string]string{
								"active-vendor":   activeVendorSource,
								"project-source":  projectSource,
								"inactive-vendor": inactiveVendorSource,
							}[entry.Source]
							if source == "" {
								return errors.New("unknown overlay source")
							}
							overlay.Replace[source] = entry.Replace
						}
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
		}
		return nil
	}

	apply, err := core.ApplyPreTestHooksWithVendorBridges(cfg, projectRoot, generatedRoot, bridges, exec)
	resp.OverlayDir = apply.OverlayDir
	resp.OverlayFile = apply.OverlayFile
	resp.GoFlags = append([]string(nil), apply.GoFlags...)
	if apply.OverlayDir != "" {
		st, statErr := os.Stat(apply.OverlayDir)
		resp.OverlayDirExists = statErr == nil && st.IsDir()
	}
	if apply.OverlayFile != "" {
		st, statErr := os.Stat(apply.OverlayFile)
		resp.OverlayFileExists = statErr == nil && !st.IsDir()
		if statErr == nil {
			resp.OverlayFileSize = st.Size()
		}
	}
	if err != nil {
		resp.ErrMsg = err.Error()
	}
	if resp.OverlayFile != "" {
		data, readErr := os.ReadFile(resp.OverlayFile)
		if readErr != nil {
			return nil, readErr
		}
		var overlay struct {
			Replace map[string]string `json:"Replace"`
		}
		if len(data) > 0 {
			if unmarshalErr := json.Unmarshal(data, &overlay); unmarshalErr != nil {
				return nil, unmarshalErr
			}
		}
		resp.OverlayReplace = overlay.Replace
	}
	return resp, nil
}
```
