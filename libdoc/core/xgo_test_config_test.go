package core

import (
	"fmt"
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

func TestLoadXgoTestConfigGoMinMax(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultXgoTestConfigName)
	body := `{"go":{"min":"1.18","max":"1.20"},"flags":["--unified"]}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadXgoTestConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Go == nil || cfg.Go.Min != "1.18" || cfg.Go.Max != "1.20" {
		t.Fatalf("Go=%#v", cfg.Go)
	}
	if len(cfg.Flags) != 1 {
		t.Fatalf("Flags=%v", cfg.Flags)
	}
}

func TestValidateXgoTestConfigGoVersion(t *testing.T) {
	if err := ValidateXgoTestConfigGoVersion(nil, "go"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateXgoTestConfigGoVersion(&XgoTestConfig{}, "go"); err != nil {
		t.Fatal(err)
	}
	// Below min.
	parsed, err := parseXgoTestConfig([]byte(`{"go":{"min":"99.0"}}`))
	if err != nil {
		t.Fatal(err)
	}
	err = ValidateXgoTestConfigGoVersion(parsed, "go")
	if err == nil || !strings.Contains(err.Error(), "< 99.0") {
		t.Fatalf("want min error, got %v", err)
	}
	// Above max.
	parsed, err = parseXgoTestConfig([]byte(`{"go":{"max":"1.0"}}`))
	if err != nil {
		t.Fatal(err)
	}
	err = ValidateXgoTestConfigGoVersion(parsed, "go")
	if err == nil || !strings.Contains(err.Error(), "> 1.0") {
		t.Fatalf("want max error, got %v", err)
	}
	// Wide range OK.
	parsed, err = parseXgoTestConfig([]byte(`{"go":{"min":"1.0","max":"99.0"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateXgoTestConfigGoVersion(parsed, "go"); err != nil {
		t.Fatal(err)
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
  "args": ["--config_file=x.ini"],
  "bypass_go_flags": true
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
	if len(got.ProgArgs) != 2 || got.ProgArgs[0] != "--" || got.ProgArgs[1] != "--config_file=x.ini" {
		t.Fatalf("ProgArgs=%v, want [-- --config_file=x.ini]", got.ProgArgs)
	}
}

func TestBuildXgoTestConfigApplySkipsArgsWithoutBypassGoFlags(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultXgoTestConfigName)
	if err := os.WriteFile(path, []byte(`{"args":["--config_dir=/tmp/config"]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := BuildXgoTestConfigApply("xgo", dir, "example.com/project", "testcase")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ProgArgs) != 0 {
		t.Fatalf("ProgArgs=%v, want no args without bypass_go_flags", got.ProgArgs)
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

func TestBuildXgoTestConfigApplyExpandsProjectRootInArgs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, DefaultXgoTestConfigName)
	body := `{
  "args": [
    "--config_dir=$PROJECT_ROOT/config",
    "--literal=$UNRELATED_ENV"
  ],
  "bypass_go_flags": true
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := BuildXgoTestConfigApply("xgo", dir, "example.com/project", "testcase")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"--", "--config_dir=" + filepath.Join(dir, "config"), "--literal=$UNRELATED_ENV"}
	if strings.Join(got.ProgArgs, "\n") != strings.Join(want, "\n") {
		t.Fatalf("ProgArgs=%q, want %q", got.ProgArgs, want)
	}
}

func TestExpandProjectRootArgsNeedsRootOnlyWhenUsed(t *testing.T) {
	got, err := expandProjectRootArgs([]string{"--literal=$PROJECT_ROOT"}, "")
	if err == nil || got != nil {
		t.Fatalf("got (%q, %v), want placeholder-root error", got, err)
	}

	got, err = expandProjectRootArgs([]string{"--literal=$UNRELATED_ENV"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "--literal=$UNRELATED_ENV" {
		t.Fatalf("got %q", got)
	}
}

func TestParsePreTestValidAndErrors(t *testing.T) {
	cfg, err := parseXgoTestConfig([]byte(`{
  "pre_test": [
    {"command": ["tool", "--flag"]},
    {"command": ["other"]}
  ]
}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.PreTest) != 2 || len(cfg.PreTest[0].Command) != 2 {
		t.Fatalf("PreTest=%#v", cfg.PreTest)
	}

	_, err = parseXgoTestConfig([]byte(`{"pre_test":[{"command":[]}]}`))
	if err == nil {
		t.Fatal("want missing-command error")
	}
	if !strings.Contains(err.Error(), "command is required") {
		t.Fatalf("err=%v", err)
	}

	_, err = parseXgoTestConfig([]byte(`{"pre_test":"nope"}`))
	if err == nil {
		t.Fatal("want wrong-type error")
	}
}

func TestApplyPreTestHooksEmpty(t *testing.T) {
	apply, err := ApplyPreTestHooks(&XgoTestConfig{}, t.TempDir(), t.TempDir(), func(string, []string) error {
		t.Fatal("executor must not run")
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if apply.HookCount != 0 || apply.OverlayDir != "" || apply.OverlayFile != "" || len(apply.GoFlags) != 0 {
		t.Fatalf("apply=%#v", apply)
	}

	apply, err = ApplyPreTestHooks(nil, t.TempDir(), t.TempDir(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if apply.HookCount != 0 {
		t.Fatalf("nil config apply=%#v", apply)
	}
}

func TestApplyPreTestHooksFlexibleMidStringAndProjectRoot(t *testing.T) {
	projectRoot := t.TempDir()
	overlayRoot := t.TempDir()
	var calls [][]string
	cfg := &XgoTestConfig{PreTest: []PreTestHook{{
		Command: []string{
			"tool",
			"--overlay=$GO_INSTRUMENT_OVERLAY_FILE",
			"--dir=$GO_INSTRUMENT_OVERLAY_DIR/extra",
			"--config=$PROJECT_ROOT/cfg",
			"--literal=$OTHER",
		},
	}}}
	apply, err := ApplyPreTestHooks(cfg, projectRoot, overlayRoot, func(workDir string, command []string) error {
		if workDir != projectRoot {
			t.Fatalf("workDir=%q want %q", workDir, projectRoot)
		}
		calls = append(calls, append([]string(nil), command...))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if apply.HookCount != 1 {
		t.Fatalf("HookCount=%d", apply.HookCount)
	}
	if apply.OverlayDir == "" || apply.OverlayFile == "" || !filepath.IsAbs(apply.OverlayDir) || !filepath.IsAbs(apply.OverlayFile) {
		t.Fatalf("overlay paths: %#v", apply)
	}
	if len(apply.GoFlags) != 0 {
		t.Fatalf("empty overlay must not add flags: %#v", apply.GoFlags)
	}
	st, err := os.Stat(apply.OverlayFile)
	if err != nil || st.Size() != 0 {
		t.Fatalf("overlay file stat: size=%v err=%v", st, err)
	}
	if len(calls) != 1 {
		t.Fatalf("calls=%#v", calls)
	}
	absRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"tool",
		"--overlay=" + apply.OverlayFile,
		"--dir=" + apply.OverlayDir + "/extra",
		"--config=" + absRoot + "/cfg",
		"--literal=$OTHER",
	}
	if strings.Join(calls[0], "\n") != strings.Join(want, "\n") {
		t.Fatalf("command=%q want %q", calls[0], want)
	}
}

func TestApplyPreTestHooksPopulatedOverlayFlag(t *testing.T) {
	projectRoot := t.TempDir()
	overlayRoot := t.TempDir()
	cfg := &XgoTestConfig{PreTest: []PreTestHook{{
		Command: []string{"writer", "--overlay-file", "$GO_INSTRUMENT_OVERLAY_FILE"},
	}}}
	apply, err := ApplyPreTestHooks(cfg, projectRoot, overlayRoot, func(_ string, command []string) error {
		return os.WriteFile(command[2], []byte(`{"Replace":{}}`), 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"-overlay=" + apply.OverlayFile}
	if len(apply.GoFlags) != 1 || apply.GoFlags[0] != want[0] {
		t.Fatalf("GoFlags=%#v want %#v", apply.GoFlags, want)
	}
}

func TestApplyPreTestHooksFailureNoGoFlags(t *testing.T) {
	projectRoot := t.TempDir()
	overlayRoot := t.TempDir()
	cfg := &XgoTestConfig{PreTest: []PreTestHook{
		{Command: []string{"ok", "$GO_INSTRUMENT_OVERLAY_FILE"}},
		{Command: []string{"fail"}},
	}}
	call := 0
	apply, err := ApplyPreTestHooks(cfg, projectRoot, overlayRoot, func(_ string, command []string) error {
		call++
		if call == 1 {
			return os.WriteFile(command[1], []byte(`{"Replace":{}}`), 0o644)
		}
		return fmt.Errorf("boom")
	})
	if err == nil {
		t.Fatal("want hook failure")
	}
	if !strings.Contains(err.Error(), "pre_test hook 2") {
		t.Fatalf("err=%v", err)
	}
	if apply.HookCount != 0 || len(apply.GoFlags) != 0 || apply.OverlayFile != "" {
		t.Fatalf("failed apply must be empty: %#v", apply)
	}
}

func TestExpandConfigArgsMidStringAllPlaceholders(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "__overlay")
	file := filepath.Join(dir, "overlay.json")
	got, err := expandConfigArgs([]string{
		"--overlay=$GO_INSTRUMENT_OVERLAY_FILE",
		"--dir=$GO_INSTRUMENT_OVERLAY_DIR/x",
		"--cfg=$PROJECT_ROOT/a",
		"$OTHER",
	}, root, dir, file)
	if err != nil {
		t.Fatal(err)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"--overlay=" + file,
		"--dir=" + dir + "/x",
		"--cfg=" + absRoot + "/a",
		"$OTHER",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("got %q want %q", got, want)
	}
}
