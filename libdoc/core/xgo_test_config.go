package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DefaultXgoTestConfigName is the project-root config file used by xgo test-explorer
// and by doctest when spawning xgo test for a generated suite.
const DefaultXgoTestConfigName = "test.config.json"

// XgoTestConfig is the subset of xgo's test.config.json that doctest applies
// when running `xgo test` against a generated suite module.
//
// Field shapes match github.com/xhd2015/xgo/cmd/xgo/test-explorer.TestConfig
// so project configs stay shared between `xgo e` and doctest.
type XgoTestConfig struct {
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
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse %s: %w", DefaultXgoTestConfigName, err)
	}
	conf := &XgoTestConfig{}
	if e, ok := m["env"]; ok && e != nil {
		em, ok := e.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s env: want object, got %T", DefaultXgoTestConfigName, e)
		}
		conf.Env = em
	}
	if e, ok := m["flags"]; ok && e != nil {
		list, err := jsonStringList(e)
		if err != nil {
			return nil, fmt.Errorf("%s flags: %w", DefaultXgoTestConfigName, err)
		}
		conf.Flags = list
	}
	if e, ok := m["args"]; ok && e != nil {
		list, err := jsonStringList(e)
		if err != nil {
			return nil, fmt.Errorf("%s args: %w", DefaultXgoTestConfigName, err)
		}
		conf.Args = list
	}
	if e, ok := m["bypass_go_flags"]; ok && e != nil {
		b, err := jsonBool(e)
		if err != nil {
			return nil, fmt.Errorf("%s bypass_go_flags: %w", DefaultXgoTestConfigName, err)
		}
		conf.BypassGoFlags = b
	}
	if e, ok := m["mock_rules"]; ok && e != nil {
		list, err := jsonMarshaledObjects(e)
		if err != nil {
			return nil, fmt.Errorf("%s mock_rules: %w", DefaultXgoTestConfigName, err)
		}
		conf.MockRules = list
	}
	return conf, nil
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
func BuildXgoTestConfigApply(goTestBin, projectModRoot, projectModPath, suiteModPath string) (XgoTestConfigApply, error) {
	var out XgoTestConfigApply
	if strings.TrimSpace(goTestBin) != "xgo" {
		return out, nil
	}
	cfgPath := FindXgoTestConfigPath(projectModRoot)
	if cfgPath == "" {
		// Still allow include-as without a config file (rules may come from CLI).
		if needIncludeAsMainModule(projectModPath, suiteModPath) {
			out.IncludeAsMainModule = projectModPath
			out.Flags = []string{"--mock-rule-include-as-main-module=" + projectModPath}
		}
		return out, nil
	}
	cfg, err := LoadXgoTestConfig(cfgPath)
	if err != nil {
		return out, err
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
	// Do not forward test.config.json "args" (e.g. --config_file=…) as go test
	// -args. Those are for xgo test-explorer / traditional binaries that parse
	// os.Args in boot.Init; doctest suite packages often never run that path and
	// flag.Parse then fails with "flag provided but not defined". Env + mock_rules
	// + flags are the portable instrumentation surface for mapping-gen suites.
	_ = cfg.Args
	_ = cfg.BypassGoFlags
	return out, nil
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

func jsonStringList(v interface{}) ([]string, error) {
	switch t := v.(type) {
	case []interface{}:
		out := make([]string, 0, len(t))
		for i, e := range t {
			s, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("index %d: want string, got %T", i, e)
			}
			out = append(out, s)
		}
		return out, nil
	case string:
		if t == "" {
			return nil, nil
		}
		return []string{t}, nil
	default:
		return nil, fmt.Errorf("want string or list, got %T", v)
	}
}

func jsonBool(v interface{}) (bool, error) {
	switch t := v.(type) {
	case bool:
		return t, nil
	default:
		return false, fmt.Errorf("want bool, got %T", v)
	}
}

// jsonMarshaledObjects re-encodes each element as a compact JSON string so
// xgo --mock-rule receives the same payload as test-explorer.
func jsonMarshaledObjects(v interface{}) ([]string, error) {
	arr, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("want list, got %T", v)
	}
	out := make([]string, 0, len(arr))
	for i, e := range arr {
		b, err := json.Marshal(e)
		if err != nil {
			return nil, fmt.Errorf("index %d: %w", i, err)
		}
		out = append(out, string(b))
	}
	return out, nil
}
