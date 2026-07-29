package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadXgoTestConfigMissing(t *testing.T) {
	cfg, err := LoadXgoTestConfig(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatal(err)
	}
	if cfg != nil {
		t.Fatalf("want nil for missing, got %#v", cfg)
	}
}

func TestLoadXgoTestConfigMockRulesAndFlags(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultXgoTestConfigName)
	body := `{
  "env": {"FOO": true, "BAR": "x"},
  "flags": ["--trap-stdlib=false", "--unified"],
  "args": ["--config_file=config-unit-test.ini"],
  "bypass_go_flags": true,
  "mock_rules": [
    {"main_module": true, "kind": "func", "action": "include"},
    {"comment": "log", "pkg": "example.com/log", "action": "include"},
    {"any": true, "action": "exclude"}
  ]
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadXgoTestConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Flags) != 2 || cfg.Flags[0] != "--trap-stdlib=false" {
		t.Fatalf("flags=%v", cfg.Flags)
	}
	if len(cfg.MockRules) != 3 {
		t.Fatalf("mock_rules len=%d", len(cfg.MockRules))
	}
	if !strings.Contains(cfg.MockRules[0], `"main_module":true`) {
		t.Fatalf("rule0=%s", cfg.MockRules[0])
	}
	if !strings.Contains(cfg.MockRules[1], "example.com/log") {
		t.Fatalf("rule1=%s", cfg.MockRules[1])
	}
	if !cfg.BypassGoFlags || len(cfg.Args) != 1 {
		t.Fatalf("args/bypass: %#v", cfg)
	}
}

func TestBuildXgoTestConfigApplySkipsNonXgo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultXgoTestConfigName)
	if err := os.WriteFile(path, []byte(`{"flags":["--trap-stdlib=false"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := BuildXgoTestConfigApply("go", dir, "example.com/proj", "testcase")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Flags) != 0 || got.ConfigPath != "" {
		t.Fatalf("want no apply for go, got %#v", got)
	}
}

func TestBuildXgoTestConfigApplyIncludeAsWhenSuiteDiffers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultXgoTestConfigName)
	body := `{
  "flags": ["--trap-stdlib=false"],
  "env": {"LOCAL_DISABLE_UID_ALLOCATOR": true},
  "mock_rules": [
    {"main_module": true, "kind": "func", "action": "include"},
    {"pkg": "example.com/log", "action": "include"},
    {"any": true, "action": "exclude"}
  ],
  "args": ["--config_file=x.ini"]
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	const proj = "example.com/pricing"
	got, err := BuildXgoTestConfigApply("xgo", dir, proj, "testcase")
	if err != nil {
		t.Fatal(err)
	}
	if got.ConfigPath != path && got.ConfigPath != filepath.Join(dir, DefaultXgoTestConfigName) {
		// abs may normalize
		if filepath.Base(got.ConfigPath) != DefaultXgoTestConfigName {
			t.Fatalf("ConfigPath=%q", got.ConfigPath)
		}
	}
	if got.IncludeAsMainModule != proj {
		t.Fatalf("IncludeAsMainModule=%q", got.IncludeAsMainModule)
	}
	joined := strings.Join(got.Flags, " ")
	if !strings.Contains(joined, "--trap-stdlib=false") {
		t.Fatalf("missing trap flag: %v", got.Flags)
	}
	if !strings.Contains(joined, "--mock-rule") {
		t.Fatalf("missing mock-rule: %v", got.Flags)
	}
	if !strings.Contains(joined, "--mock-rule-include-as-main-module="+proj) {
		t.Fatalf("missing include-as: %v", got.Flags)
	}
	// mock rules present as separate argv pairs
	mockN := 0
	for i := 0; i < len(got.Flags); i++ {
		if got.Flags[i] == "--mock-rule" {
			mockN++
		}
	}
	if mockN != 3 {
		t.Fatalf("mock-rule count=%d flags=%v", mockN, got.Flags)
	}
	if len(got.Env) != 1 || !strings.HasPrefix(got.Env[0], "LOCAL_DISABLE_UID_ALLOCATOR=") {
		t.Fatalf("env=%v", got.Env)
	}
	// ProgArgs intentionally omitted for doctest suites (see BuildXgoTestConfigApply).
	if len(got.ProgArgs) != 0 {
		t.Fatalf("ProgArgs should be empty, got %v", got.ProgArgs)
	}
}

func TestBuildXgoTestConfigApplyNoIncludeWhenSameModule(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultXgoTestConfigName)
	if err := os.WriteFile(path, []byte(`{"mock_rules":[{"any":true,"action":"exclude"}]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	const proj = "example.com/pricing"
	got, err := BuildXgoTestConfigApply("xgo", dir, proj, proj)
	if err != nil {
		t.Fatal(err)
	}
	if got.IncludeAsMainModule != "" {
		t.Fatalf("want no include-as for same module, got %q", got.IncludeAsMainModule)
	}
	for _, f := range got.Flags {
		if strings.HasPrefix(f, "--mock-rule-include-as-main-module") {
			t.Fatalf("unexpected include-as flag: %v", got.Flags)
		}
	}
}

func TestNeedIncludeAsMainModule(t *testing.T) {
	if !needIncludeAsMainModule("a", "b") {
		t.Fatal("diff should include")
	}
	if needIncludeAsMainModule("a", "a") {
		t.Fatal("same should not")
	}
	if !needIncludeAsMainModule("a", "") {
		t.Fatal("empty suite should include")
	}
	if needIncludeAsMainModule("", "testcase") {
		t.Fatal("empty project should not")
	}
}
