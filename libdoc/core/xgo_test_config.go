package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/xgo/support/testconfig"
)

// DefaultXgoTestConfigName is the project-root config file used by xgo test-explorer
// and by doctest when spawning xgo test for a generated suite.
const DefaultXgoTestConfigName = testconfig.DefaultFileName

// XgoTestConfig is the subset of xgo's test.config.json that doctest applies
// when running `xgo test` against a generated suite module.
//
// Shared schema fields (go min/max, env, flags, args, mock_rules, …) are parsed
// via github.com/xhd2015/xgo/support/testconfig. PreTest is doctest-specific.
type XgoTestConfig struct {
	// Go is optional go.min / go.max from test.config.json (shared with xgo).
	Go *testconfig.GoConfig `json:"go"`
	// PreTest are generic commands run before the generated suite is compiled.
	// They may use the exact overlay placeholders documented by PreTestHook.
	PreTest []PreTestHook `json:"pre_test"`
	// Env is KEY -> any JSON value (stringified when exported as KEY=value).
	Env map[string]interface{} `json:"env"`
	// Flags are xgo/go build-test flags (e.g. --trap-stdlib=false).
	Flags []string `json:"flags"`
	// Args are test-binary program args (after -args).
	Args []string `json:"args"`
	// BypassGoFlags mirrors xgo explorer: when true with Args, still append -args.
	BypassGoFlags bool `json:"bypass_go_flags"`
	// MockRules are raw JSON objects (re-marshaled for --mock-rule).
	MockRules []string `json:"-"`
}

// PreTestHook is one ordered command declared by test.config.json. Command is
// an argv list, not a shell expression.
type PreTestHook struct {
	Command []string `json:"command"`
}

// PreTestHookExecutor runs a fully expanded hook from workDir.
type PreTestHookExecutor func(workDir string, command []string) error

// PreTestHookApply is the generic contribution from pre-test hooks. GoFlags
// is either empty or contains one standard Go -overlay argument.
type PreTestHookApply struct {
	OverlayDir  string
	OverlayFile string
	GoFlags     []string
	// HookCount is the number of hooks that ran successfully (0 when none).
	HookCount int
}

// VendorBridgeMapping identifies a vendored module that needed a synthetic
// go.mod for the current generated workspace (xgo-style). BridgeRoot is the
// absolute path of the placeholder go.mod file under vendor-gomod-overlay;
// OriginalVendorRoot is the project vendor/<mod> directory (replace target).
// Callers must not infer mappings by scanning gen cache directories.
type VendorBridgeMapping struct {
	ModulePath         string
	OriginalVendorRoot string
	BridgeRoot         string // absolute placeholder go.mod path
}

const (
	goInstrumentOverlayDirPlaceholder  = "$GO_INSTRUMENT_OVERLAY_DIR"
	goInstrumentOverlayFilePlaceholder = "$GO_INSTRUMENT_OVERLAY_FILE"
	// projectRootPlaceholder is deliberately expanded by doctest rather than
	// the process environment: test.config.json must remain relocatable.
	projectRootPlaceholder = "$PROJECT_ROOT"
)

// knownConfigPlaceholders is the fixed expand order for substring replacement.
// Unknown $TOKEN values are never touched.
var knownConfigPlaceholders = []string{
	goInstrumentOverlayDirPlaceholder,
	goInstrumentOverlayFilePlaceholder,
	projectRootPlaceholder,
}

// LoadXgoTestConfig reads absPath (typically <projectModRoot>/test.config.json).
// Missing file returns (nil, nil). Invalid JSON returns an error.
func LoadXgoTestConfig(absPath string) (*XgoTestConfig, error) {
	if absPath == "" {
		return nil, nil
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return &XgoTestConfig{}, nil
	}
	return parseXgoTestConfig(data)
}

// FindXgoTestConfigPath returns the absolute path of test.config.json under
// modRoot when present. Empty string means not found.
func FindXgoTestConfigPath(modRoot string) string {
	if modRoot == "" {
		return ""
	}
	p := filepath.Join(modRoot, DefaultXgoTestConfigName)
	st, err := os.Stat(p)
	if err != nil || st.IsDir() {
		return ""
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

func parseXgoTestConfig(data []byte) (*XgoTestConfig, error) {
	shared, err := testconfig.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", DefaultXgoTestConfigName, err)
	}
	conf := &XgoTestConfig{
		Go:            shared.Go,
		Env:           shared.Env,
		Flags:         shared.Flags,
		Args:          shared.Args,
		BypassGoFlags: shared.BypassGoFlags,
		MockRules:     shared.MockRules,
	}
	// pre_test is doctest-only; parse from the raw map after shared schema.
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", DefaultXgoTestConfigName, err)
	}
	if e, ok := m["pre_test"]; ok && e != nil {
		hooks, err := jsonPreTestHooks(e)
		if err != nil {
			return nil, fmt.Errorf("%s pre_test: %w", DefaultXgoTestConfigName, err)
		}
		conf.PreTest = hooks
	}
	return conf, nil
}

// ValidateXgoTestConfigGoVersion enforces test.config.json go.min/go.max against
// the host go toolchain (PATH "go"), matching xgo test-explorer. No-op when
// constraints are empty. goBinary defaults to "go".
func ValidateXgoTestConfigGoVersion(cfg *XgoTestConfig, goBinary string) error {
	if cfg == nil {
		return nil
	}
	return testconfig.ValidateGoConstraint(cfg.Go, goBinary)
}

// ApplyPreTestHooks expands unified config placeholders (substring Contains /
// ReplaceAll), executes hooks in declaration order, and returns a Go overlay
// flag only when a hook populated the driver-owned file. overlayRoot is a
// durable, source-project mapping-cache root controlled by the caller.
//
// Placeholder vocabulary: $PROJECT_ROOT, $GO_INSTRUMENT_OVERLAY_DIR,
// $GO_INSTRUMENT_OVERLAY_FILE. Unknown $TOKEN values are left untouched.
func ApplyPreTestHooks(config *XgoTestConfig, projectRoot string, overlayRoot string, execute PreTestHookExecutor) (PreTestHookApply, error) {
	return ApplyPreTestHooksWithVendorBridges(config, projectRoot, overlayRoot, nil, execute)
}

// ApplyPreTestHooksWithVendorBridges is ApplyPreTestHooks with explicit
// current-run vendor bridge metadata (xgo-style). After all hooks succeed, it
// merges phantom vendor go.mod mappings into the shared overlay Replace map:
// project vendor/<mod>/go.mod → placeholder under vendor-gomod-overlay.
// Package overlay keys stay on project vendor/; replacement values are never
// rewritten.
func ApplyPreTestHooksWithVendorBridges(config *XgoTestConfig, projectRoot string, overlayRoot string, bridges []VendorBridgeMapping, execute PreTestHookExecutor) (PreTestHookApply, error) {
	return ApplyPreTestHooksWithUserOverlay(config, projectRoot, overlayRoot, "", bridges, execute)
}

// ApplyPreTestHooksWithUserOverlay is ApplyPreTestHooksWithVendorBridges with
// an optional user overlay seed. When userOverlay is non-empty, the driver
// __overlay/overlay.json is opened and seeded with the user's Replace map
// before hooks run (even if hooks omit $GO_INSTRUMENT_OVERLAY_* placeholders),
// so seed alone can still emit a -overlay= flag after vendor-bridge merge.
// Later layers win the same Replace disk-path key: hook writes, then vendor
// bridges. Empty userOverlay matches ApplyPreTestHooksWithVendorBridges.
func ApplyPreTestHooksWithUserOverlay(config *XgoTestConfig, projectRoot string, overlayRoot string, userOverlay string, bridges []VendorBridgeMapping, execute PreTestHookExecutor) (PreTestHookApply, error) {
	var out PreTestHookApply
	userOverlay = strings.TrimSpace(userOverlay)
	if config == nil || len(config.PreTest) == 0 {
		return out, nil
	}
	if strings.TrimSpace(projectRoot) == "" {
		return out, fmt.Errorf("pre_test: project root is required")
	}
	if strings.TrimSpace(overlayRoot) == "" {
		return out, fmt.Errorf("pre_test: overlay root is required")
	}
	if execute == nil {
		return out, fmt.Errorf("pre_test: hook executor is required")
	}

	needDir, needFile, needProjectRoot := false, false, false
	for i, hook := range config.PreTest {
		if len(hook.Command) == 0 {
			return out, fmt.Errorf("pre_test hook %d: command is required", i+1)
		}
		for _, arg := range hook.Command {
			if strings.Contains(arg, goInstrumentOverlayDirPlaceholder) {
				needDir = true
			}
			if strings.Contains(arg, goInstrumentOverlayFilePlaceholder) {
				needFile = true
			}
			if strings.Contains(arg, projectRootPlaceholder) {
				needProjectRoot = true
			}
		}
	}
	// User seed requires a driver overlay file even when hooks omit placeholders.
	if userOverlay != "" {
		needFile = true
	}

	if needDir || needFile {
		dir, err := filepath.Abs(filepath.Join(overlayRoot, "__overlay"))
		if err != nil {
			return out, fmt.Errorf("pre_test: resolve overlay directory: %w", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return out, fmt.Errorf("pre_test: create overlay directory: %w", err)
		}
		if needDir {
			out.OverlayDir = dir
		}
		if needFile {
			out.OverlayFile = filepath.Join(dir, "overlay.json")
			if userOverlay != "" {
				if err := writeSeededOverlayFile(out.OverlayFile, userOverlay); err != nil {
					return out, err
				}
			} else {
				file, err := os.OpenFile(out.OverlayFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
				if err != nil {
					return out, fmt.Errorf("pre_test: create overlay file: %w", err)
				}
				if err := file.Close(); err != nil {
					return out, fmt.Errorf("pre_test: close overlay file: %w", err)
				}
			}
		}
	}

	vars := map[string]string{}
	if needProjectRoot {
		absRoot, err := filepath.Abs(projectRoot)
		if err != nil {
			return out, fmt.Errorf("pre_test: resolve %s: %w", projectRootPlaceholder, err)
		}
		vars[projectRootPlaceholder] = absRoot
	}
	if out.OverlayDir != "" {
		vars[goInstrumentOverlayDirPlaceholder] = out.OverlayDir
	}
	if out.OverlayFile != "" {
		vars[goInstrumentOverlayFilePlaceholder] = out.OverlayFile
	}

	for i, hook := range config.PreTest {
		command := expandConfigPlaceholders(hook.Command, vars)
		if err := execute(projectRoot, command); err != nil {
			return PreTestHookApply{}, fmt.Errorf("pre_test hook %d (%s): %w", i+1, strings.Join(command, " "), err)
		}
	}

	if out.OverlayFile != "" {
		st, err := os.Stat(out.OverlayFile)
		if err != nil {
			return PreTestHookApply{}, fmt.Errorf("pre_test: inspect overlay file: %w", err)
		}
		if st.Size() > 0 {
			if err := normalizeOverlayVendorBridgeSources(out.OverlayFile, bridges); err != nil {
				return PreTestHookApply{}, err
			}
			out.GoFlags = []string{"-overlay=" + out.OverlayFile}
		}
	}
	out.HookCount = len(config.PreTest)
	return out, nil
}

// MaterializeUserVendorOverlay merges a user overlay seed with a vendor-gomod
// overlay JSON when there is no pre_test driver file path. Vendor Replace keys
// overwrite the same disk-path keys from the user seed. Writes the merged
// result under destRoot/__overlay/overlay.json and returns at most one
// -overlay= GoFlags entry when the merged Replace map is non-empty.
// Empty user and vendor paths (or empty Replace) yield a zero apply with no error.
func MaterializeUserVendorOverlay(userOverlay, vendorGomodOverlay, destRoot string) (PreTestHookApply, error) {
	var out PreTestHookApply
	userOverlay = strings.TrimSpace(userOverlay)
	vendorGomodOverlay = strings.TrimSpace(vendorGomodOverlay)
	if userOverlay == "" && vendorGomodOverlay == "" {
		return out, nil
	}
	if strings.TrimSpace(destRoot) == "" {
		return out, fmt.Errorf("materialize overlay: dest root is required")
	}

	merged := make(map[string]string)
	if userOverlay != "" {
		userReplace, err := readOverlayReplaceMap(userOverlay)
		if err != nil {
			return out, fmt.Errorf("materialize overlay: read user overlay: %w", err)
		}
		for k, v := range userReplace {
			if k == "" {
				continue
			}
			merged[k] = v
		}
	}
	if vendorGomodOverlay != "" {
		vendorReplace, err := readOverlayReplaceMap(vendorGomodOverlay)
		if err != nil {
			return out, fmt.Errorf("materialize overlay: read vendor overlay: %w", err)
		}
		for k, v := range vendorReplace {
			if k == "" {
				continue
			}
			// Vendor later-wins on the same Replace key.
			merged[k] = v
		}
	}
	if len(merged) == 0 {
		return out, nil
	}

	dir, err := filepath.Abs(filepath.Join(destRoot, "__overlay"))
	if err != nil {
		return out, fmt.Errorf("materialize overlay: resolve dest: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return out, fmt.Errorf("materialize overlay: create dest: %w", err)
	}
	outPath := filepath.Join(dir, "overlay.json")
	payload, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{Replace: merged})
	if err != nil {
		return out, fmt.Errorf("materialize overlay: encode: %w", err)
	}
	if err := os.WriteFile(outPath, payload, 0o644); err != nil {
		return out, fmt.Errorf("materialize overlay: write: %w", err)
	}
	abs, err := filepath.Abs(outPath)
	if err != nil {
		abs = outPath
	}
	out.OverlayDir = dir
	out.OverlayFile = abs
	out.GoFlags = []string{"-overlay=" + abs}
	return out, nil
}

// writeSeededOverlayFile writes driver overlay JSON seeded from userOverlay's
// Replace map (empty map when the user file has no Replace entries).
func writeSeededOverlayFile(destPath, userOverlay string) error {
	replace, err := readOverlayReplaceMap(userOverlay)
	if err != nil {
		return fmt.Errorf("pre_test: read user overlay seed: %w", err)
	}
	if replace == nil {
		replace = map[string]string{}
	}
	payload, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{Replace: replace})
	if err != nil {
		return fmt.Errorf("pre_test: encode user overlay seed: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("pre_test: create overlay directory: %w", err)
	}
	if err := os.WriteFile(destPath, payload, 0o644); err != nil {
		return fmt.Errorf("pre_test: write user overlay seed: %w", err)
	}
	return nil
}

// readOverlayReplaceMap loads the Replace map from a standard Go overlay JSON.
// Missing file is an error; empty/whitespace file yields an empty map.
func readOverlayReplaceMap(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return map[string]string{}, nil
	}
	var payload struct {
		Replace map[string]string `json:"Replace"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	if payload.Replace == nil {
		return map[string]string{}, nil
	}
	return payload.Replace, nil
}

// normalizeOverlayVendorBridgeSources merges xgo-style phantom vendor go.mod
// mappings into a standard Go overlay Replace map. Packages stay under project
// vendor/; only missing go.mod files are overlaid. Does not rewrite other
// overlay source keys.
func normalizeOverlayVendorBridgeSources(overlayFile string, bridges []VendorBridgeMapping) error {
	if len(bridges) == 0 {
		return nil
	}
	data, err := os.ReadFile(overlayFile)
	if err != nil {
		return fmt.Errorf("pre_test: read overlay file: %w", err)
	}
	var overlay map[string]json.RawMessage
	if len(strings.TrimSpace(string(data))) == 0 {
		overlay = map[string]json.RawMessage{}
	} else if err := json.Unmarshal(data, &overlay); err != nil {
		return fmt.Errorf("pre_test: parse overlay file: %w", err)
	}
	replace := map[string]string{}
	if replaceRaw, ok := overlay["Replace"]; ok {
		if err := json.Unmarshal(replaceRaw, &replace); err != nil {
			return fmt.Errorf("pre_test: parse overlay Replace: %w", err)
		}
	}
	changed := false
	for _, bridge := range bridges {
		src, dst, ok := vendorGoModOverlayPair(bridge)
		if !ok {
			continue
		}
		if replace[src] != dst {
			replace[src] = dst
			changed = true
		}
	}
	if !changed && len(replace) == 0 {
		return nil
	}
	if !changed {
		// Still rewrite if file was empty but replace already complete via other means.
		if len(data) > 0 {
			return nil
		}
	}
	replaceData, err := json.Marshal(replace)
	if err != nil {
		return fmt.Errorf("pre_test: encode overlay Replace: %w", err)
	}
	if overlay == nil {
		overlay = map[string]json.RawMessage{}
	}
	overlay["Replace"] = replaceData
	out, err := json.Marshal(overlay)
	if err != nil {
		return fmt.Errorf("pre_test: encode overlay file: %w", err)
	}
	if err := os.WriteFile(overlayFile, out, 0o644); err != nil {
		return fmt.Errorf("pre_test: write overlay file: %w", err)
	}
	return nil
}

func vendorGoModOverlayPair(bridge VendorBridgeMapping) (src, dst string, ok bool) {
	if strings.TrimSpace(bridge.OriginalVendorRoot) == "" || strings.TrimSpace(bridge.BridgeRoot) == "" {
		return "", "", false
	}
	srcPath := filepath.Join(bridge.OriginalVendorRoot, "go.mod")
	srcAbs, err := filepath.Abs(srcPath)
	if err != nil {
		srcAbs = srcPath
	}
	dstAbs, err := filepath.Abs(bridge.BridgeRoot)
	if err != nil {
		dstAbs = bridge.BridgeRoot
	}
	return srcAbs, dstAbs, true
}

// XgoTestConfigApply describes how project test.config.json is applied when
// suite module path differs from the project (mapping-gen / testcase module).
type XgoTestConfigApply struct {
	// ConfigPath is the absolute project config path used (empty if none).
	ConfigPath string
	// Flags are extra xgo argv fragments after "test" (mock-rule, include-as, flags).
	Flags []string
	// Env is KEY=value overrides for the child xgo/go process.
	Env []string
	// ProgArgs are test-binary args (after packages / -args).
	ProgArgs []string
	// IncludeAsMainModule is the project module path when suite ≠ project.
	IncludeAsMainModule string
}

// BuildXgoTestConfigApply loads the project config and builds argv/env for an
// xgo test invocation.
//
// projectModRoot / projectModPath identify the consumer project (where
// test.config.json lives). suiteModPath is the module path of the package under
// test (gen module is typically "testcase"). When suiteModPath is non-empty and
// differs from projectModPath, --mock-rule-include-as-main-module is set so
// mock_rules with main_module:true apply to project packages under replace.
//
// No-op when goTestBin is not "xgo" or when no config exists.
// Prefer ApplyLoadedXgoTestConfig when the caller already loaded the config
// (shared path with pre_test application).
func BuildXgoTestConfigApply(goTestBin, projectModRoot, projectModPath, suiteModPath string) (XgoTestConfigApply, error) {
	if strings.TrimSpace(goTestBin) != "xgo" {
		return XgoTestConfigApply{}, nil
	}
	cfgPath := FindXgoTestConfigPath(projectModRoot)
	var cfg *XgoTestConfig
	if cfgPath != "" {
		var err error
		cfg, err = LoadXgoTestConfig(cfgPath)
		if err != nil {
			return XgoTestConfigApply{}, err
		}
	}
	return ApplyLoadedXgoTestConfig(goTestBin, cfgPath, cfg, projectModRoot, projectModPath, suiteModPath)
}

// ApplyLoadedXgoTestConfig builds argv/env from an already-loaded config.
// cfgPath may be empty when no file exists (include-as-only path). cfg may be
// nil. No-op when goTestBin is not "xgo".
func ApplyLoadedXgoTestConfig(goTestBin, cfgPath string, cfg *XgoTestConfig, projectModRoot, projectModPath, suiteModPath string) (XgoTestConfigApply, error) {
	var out XgoTestConfigApply
	if strings.TrimSpace(goTestBin) != "xgo" {
		return out, nil
	}
	if cfgPath == "" {
		// Still allow include-as without a config file (rules may come from CLI).
		if needIncludeAsMainModule(projectModPath, suiteModPath) {
			out.IncludeAsMainModule = projectModPath
			out.Flags = []string{"--mock-rule-include-as-main-module=" + projectModPath}
		}
		return out, nil
	}
	out.ConfigPath = cfgPath
	if cfg == nil {
		cfg = &XgoTestConfig{}
	}

	// Config flags first (trap-stdlib, unified, …), then mock rules (higher
	// priority in xgo merge semantics), then include-as.
	out.Flags = append(out.Flags, cfg.Flags...)
	for _, rule := range cfg.MockRules {
		if rule == "" {
			continue
		}
		out.Flags = append(out.Flags, "--mock-rule", rule)
	}
	if needIncludeAsMainModule(projectModPath, suiteModPath) {
		out.IncludeAsMainModule = projectModPath
		out.Flags = append(out.Flags, "--mock-rule-include-as-main-module="+projectModPath)
	}
	out.Env = xgoConfigEnvPairs(cfg.Env)
	// Mirror xgo test-explorer's bypass-go-flags behavior: the leading "--" is
	// passed after go test -args. Legacy boot code can still inspect os.Args,
	// while Go's test flag parser stops before treating these as test flags.
	if cfg.BypassGoFlags && len(cfg.Args) > 0 {
		args, err := expandConfigArgs(cfg.Args, projectModRoot, "", "")
		if err != nil {
			return out, err
		}
		out.ProgArgs = make([]string, 0, len(cfg.Args)+1)
		out.ProgArgs = append(out.ProgArgs, "--")
		out.ProgArgs = append(out.ProgArgs, args...)
	}
	return out, nil
}

// expandConfigPlaceholders applies known placeholder → value replacements as
// substrings (ReplaceAll). Unknown $TOKEN text is left unchanged. vars keys
// should be the full placeholder strings (e.g. "$PROJECT_ROOT").
func expandConfigPlaceholders(args []string, vars map[string]string) []string {
	if len(args) == 0 {
		return nil
	}
	out := make([]string, len(args))
	for i, arg := range args {
		s := arg
		for _, ph := range knownConfigPlaceholders {
			val, ok := vars[ph]
			if !ok || val == "" {
				continue
			}
			if strings.Contains(s, ph) {
				s = strings.ReplaceAll(s, ph, val)
			}
		}
		out[i] = s
	}
	return out
}

// expandConfigArgs builds the vars map from projectRoot / overlay paths and
// expands known placeholders. projectRoot is required only when $PROJECT_ROOT
// appears. Overlay placeholders expand only when the corresponding path is set
// (pre_test allocation path); otherwise they remain literal.
func expandConfigArgs(args []string, projectRoot, overlayDir, overlayFile string) ([]string, error) {
	if len(args) == 0 {
		return nil, nil
	}
	vars := map[string]string{}
	needRoot := false
	for _, arg := range args {
		if strings.Contains(arg, projectRootPlaceholder) {
			needRoot = true
			break
		}
	}
	if needRoot {
		if strings.TrimSpace(projectRoot) == "" {
			return nil, fmt.Errorf("%s requires a project module root", projectRootPlaceholder)
		}
		abs, err := filepath.Abs(projectRoot)
		if err != nil {
			return nil, fmt.Errorf("resolve %s: %w", projectRootPlaceholder, err)
		}
		vars[projectRootPlaceholder] = abs
	}
	if overlayDir != "" {
		vars[goInstrumentOverlayDirPlaceholder] = overlayDir
	}
	if overlayFile != "" {
		vars[goInstrumentOverlayFilePlaceholder] = overlayFile
	}
	if len(vars) == 0 {
		return append([]string(nil), args...), nil
	}
	return expandConfigPlaceholders(args, vars), nil
}

// expandProjectRootArgs replaces doctest's explicit $PROJECT_ROOT placeholder
// with the absolute source module root. It intentionally does not expand any
// other $NAME text from the host process environment.
func expandProjectRootArgs(args []string, projectModRoot string) ([]string, error) {
	return expandConfigArgs(args, projectModRoot, "", "")
}

func needIncludeAsMainModule(projectModPath, suiteModPath string) bool {
	projectModPath = strings.TrimSpace(projectModPath)
	suiteModPath = strings.TrimSpace(suiteModPath)
	if projectModPath == "" {
		return false
	}
	if suiteModPath == "" {
		// Unknown suite module: still include project when we have a path so
		// mapping-gen (module testcase) works without resolving gen go.mod.
		return true
	}
	return suiteModPath != projectModPath
}

func xgoConfigEnvPairs(env map[string]interface{}) []string {
	if len(env) == 0 {
		return nil
	}
	// Stable order for tests / display.
	keys := make([]string, 0, len(env))
	for k := range env {
		if k == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+"="+fmt.Sprint(env[k]))
	}
	return out
}

func jsonPreTestHooks(v interface{}) ([]PreTestHook, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var hooks []PreTestHook
	if err := json.Unmarshal(data, &hooks); err != nil {
		return nil, err
	}
	for i, hook := range hooks {
		if len(hook.Command) == 0 {
			return nil, fmt.Errorf("index %d: command is required", i)
		}
	}
	return hooks, nil
}
