package core

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestKindAShimImport(t *testing.T) {
	t.Parallel()
	got, ok := KindAShimImport("testcase/tests/http/internal/post-succeeds")
	if !ok {
		t.Fatal("expected ok")
	}
	want := "testcase/tests/http/__doctest_internal_shim/post-succeeds"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if _, ok := KindAShimImport("testcase/tests/http/external/blocked"); ok {
		t.Fatal("expected no shim for path without internal")
	}
}

func TestRewriteKindALeafImports(t *testing.T) {
	t.Parallel()
	in := []string{
		"testcase/t/a",
		"testcase/t/http/internal/x",
	}
	out := RewriteKindALeafImports(in)
	if out[0] != in[0] {
		t.Fatalf("unchanged: %q", out[0])
	}
	if out[1] != "testcase/t/http/__doctest_internal_shim/x" {
		t.Fatalf("shimmed: %q", out[1])
	}
}

func TestProductInternalImport(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		imp           string
		parentModPath string
		wantOK        bool
		wantMod       string
		wantFull      string
	}{
		{
			// P1: parent-module internal is Kind B product-internal (expose), not skipped.
			name:          "parent_internal_kind_b",
			imp:           "example.com/app/internal/greet",
			parentModPath: "example.com/app",
			wantOK:        true,
			wantMod:       "example.com/app",
			wantFull:      "example.com/app/internal/greet",
		},
		{
			name:          "external_product_internal",
			imp:           "example.com/app/internal/greet",
			parentModPath: "example.com/runner",
			wantOK:        true,
			wantMod:       "example.com/app",
			wantFull:      "example.com/app/internal/greet",
		},
		{
			name:          "testcase_kind_a_not_product",
			imp:           "testcase/t/http/internal/x",
			parentModPath: "example.com/app",
			wantOK:        false,
		},
		{
			name:          "testcase_empty_parent",
			imp:           "testcase/x/internal/y",
			parentModPath: "",
			wantOK:        false,
		},
		{
			name:          "no_internal_segment",
			imp:           "example.com/app/pkg",
			parentModPath: "example.com/app",
			wantOK:        false,
		},
		{
			name:          "nested_parent_internal",
			imp:           "example.com/app/internal/foo/bar",
			parentModPath: "example.com/app",
			wantOK:        true,
			wantMod:       "example.com/app",
			wantFull:      "example.com/app/internal/foo/bar",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mod, full, ok := productInternalImport(tt.imp, tt.parentModPath)
			if ok != tt.wantOK {
				t.Fatalf("ok=%v want %v (mod=%q full=%q)", ok, tt.wantOK, mod, full)
			}
			if !tt.wantOK {
				return
			}
			if mod != tt.wantMod || full != tt.wantFull {
				t.Fatalf("got mod=%q full=%q want mod=%q full=%q", mod, full, tt.wantMod, tt.wantFull)
			}
		})
	}
}

// TestGenerateExposeSource_exportSurface asserts the Kind B facade re-exports
// exported funcs, types, vars, and consts, and never re-exports unexported
// symbols. Package name must match the internal package (call sites stay greet.X).
//
// P1 RED until generateExposeSource emits var/const aliases.
func TestGenerateExposeSource_exportSurface(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	src := `package greet

type Info struct{ N int }

var DefaultName = "world"

const MaxN = 3

func Hello(name string) string { return "hello " + name }

func secret() string { return "x" }
`
	if err := os.WriteFile(filepath.Join(dir, "greet.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	internalImp := "example.com/app/internal/greet"
	body, err := generateExposeSource(internalImp, dir)
	if err != nil {
		t.Fatalf("generateExposeSource: %v", err)
	}

	// Package name matches internal so callers keep greet.X.
	if !strings.Contains(body, "package greet\n") {
		t.Fatalf("package name must be greet; body:\n%s", body)
	}

	// Import of the real internal path (any alias).
	if !strings.Contains(body, strconv.Quote(internalImp)) {
		t.Fatalf("missing import of %q; body:\n%s", internalImp, body)
	}

	// Type alias re-export (existing behavior — keep green).
	if !strings.Contains(body, "type Info =") || !strings.Contains(body, ".Info") {
		t.Fatalf("missing type Info alias re-export; body:\n%s", body)
	}

	// Func wrapper (existing — keep green).
	if !strings.Contains(body, "func Hello(") {
		t.Fatalf("missing func Hello wrapper; body:\n%s", body)
	}

	// Var re-export — P1 NEW (RED today: funcs/types only).
	if !hasVarReexport(body, "DefaultName") {
		t.Fatalf("missing var DefaultName re-export (e.g. var DefaultName = <alias>.DefaultName); body:\n%s", body)
	}

	// Const re-export (const MaxN = <alias>.MaxN or var fallback) — P1 NEW.
	if !hasConstOrVarReexport(body, "MaxN") {
		t.Fatalf("missing const/var MaxN re-export; body:\n%s", body)
	}

	// Unexported must not appear as a re-export.
	if strings.Contains(body, "func secret") || strings.Contains(body, ".secret") {
		t.Fatalf("must not re-export unexported secret; body:\n%s", body)
	}
}

func TestGenerateExposeSource_emptyPackage(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	// Package with only unexported symbols still needs a valid facade package clause.
	src := `package empty

func hidden() {}
`
	if err := os.WriteFile(filepath.Join(dir, "empty.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	body, err := generateExposeSource("example.com/m/internal/empty", dir)
	if err != nil {
		t.Fatalf("generateExposeSource: %v", err)
	}
	if !strings.Contains(body, "package empty") {
		t.Fatalf("want package empty; body:\n%s", body)
	}
	if strings.Contains(body, "func hidden") {
		t.Fatalf("must not export hidden; body:\n%s", body)
	}
}

// TestGenerateExposeSource_externalSigTypes imports product packages used in
// exported signatures (Kind B compile fix for undefined: model).
func TestGenerateExposeSource_externalSigTypes(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	src := `package rules

import "example.com/app/model"

func FixIgnore(project model.Project, dryRun bool) (model.FixResult, error) {
	_ = dryRun
	return model.FixResult{OK: project.Root != ""}, nil
}
`
	if err := os.WriteFile(filepath.Join(dir, "rules.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	body, err := generateExposeSource("example.com/app/internal/rules", dir)
	if err != nil {
		t.Fatalf("generateExposeSource: %v", err)
	}
	if !strings.Contains(body, strconv.Quote("example.com/app/model")) {
		t.Fatalf("missing import of model package; body:\n%s", body)
	}
	if !strings.Contains(body, "func FixIgnore(project model.Project") {
		t.Fatalf("missing FixIgnore with model.Project; body:\n%s", body)
	}
	if !strings.Contains(body, "srcpkg.FixIgnore") {
		t.Fatalf("missing forward to srcpkg; body:\n%s", body)
	}
}

func TestDefaultImportAlias(t *testing.T) {
	t.Parallel()
	tests := []struct {
		path, want string
	}{
		{"example.com/app/model", "model"},
		{"example.com/foo/bar/v2", "bar"},
		{"gopkg.in/yaml.v3", "yaml"},
		{"gopkg.in/check.v1", "check"},
		{"context", "context"},
		{"net/http", "http"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			if got := defaultImportAlias(tt.path); got != tt.want {
				t.Fatalf("defaultImportAlias(%q)=%q want %q", tt.path, got, tt.want)
			}
		})
	}
}

// TestGenerateExposeSource_v2PackageName uses the package clause (ext), not
// the last path element (v2), for unaliased …/v2 imports.
func TestGenerateExposeSource_v2PackageName(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/app\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	extDir := filepath.Join(root, "ext", "v2")
	if err := os.MkdirAll(extDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(extDir, "ext.go"), []byte("package ext\n\ntype T struct{}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "internal", "rules")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	src := `package rules

import "example.com/app/ext/v2"

func Use(t ext.T) ext.T { return t }
`
	if err := os.WriteFile(filepath.Join(dir, "rules.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	body, err := generateExposeSource("example.com/app/internal/rules", dir)
	if err != nil {
		t.Fatalf("generateExposeSource: %v", err)
	}
	if !strings.Contains(body, strconv.Quote("example.com/app/ext/v2")) {
		t.Fatalf("missing import of ext/v2; body:\n%s", body)
	}
	if !strings.Contains(body, "func Use(t ext.T)") {
		t.Fatalf("want ext.T from package clause, not last path element; body:\n%s", body)
	}
	if strings.Contains(body, "v2.T") {
		t.Fatalf("must not use path last-element v2 as selector; body:\n%s", body)
	}
}

// TestGenerateExposeSource_gopkgInPackageName uses the yaml.v3 → yaml heuristic
// when the imported package is not a sibling on disk.
func TestGenerateExposeSource_gopkgInPackageName(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	src := `package rules

import "gopkg.in/yaml.v3"

func Dump(n yaml.Node) yaml.Node { return n }
`
	if err := os.WriteFile(filepath.Join(dir, "rules.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}
	body, err := generateExposeSource("example.com/app/internal/rules", dir)
	if err != nil {
		t.Fatalf("generateExposeSource: %v", err)
	}
	if !strings.Contains(body, strconv.Quote("gopkg.in/yaml.v3")) {
		t.Fatalf("missing yaml.v3 import; body:\n%s", body)
	}
	if !strings.Contains(body, "func Dump(n yaml.Node)") {
		t.Fatalf("want yaml.Node selector; body:\n%s", body)
	}
	if strings.Contains(body, "yaml.v3.") {
		t.Fatalf("must not emit invalid selector yaml.v3; body:\n%s", body)
	}
}

func TestDoctestInternalExposeDir(t *testing.T) {
	t.Parallel()
	if DoctestInternalExposeDir != "__doctest_internal_expose" {
		t.Fatalf("DoctestInternalExposeDir = %q", DoctestInternalExposeDir)
	}
	if DoctestInternalShimDir != "__doctest_internal_shim" {
		t.Fatalf("DoctestInternalShimDir = %q", DoctestInternalShimDir)
	}
}

// TestCleanupKindBMaterialized removes session-scoped expose.go under product
// tree and prunes empty __doctest_internal_expose dirs.
func TestCleanupKindBMaterialized(t *testing.T) {
	t.Parallel()
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(filepath.Dir(virt), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(virt, []byte("package greet\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := recordKindBMaterialized(genRoot, virt); err != nil {
		t.Fatal(err)
	}
	if err := CleanupKindBMaterialized(genRoot); err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("expose.go still present: %v", err)
	}
	if _, err := os.Stat(filepath.Join(product, DoctestInternalExposeDir)); !os.IsNotExist(err) {
		t.Fatalf("expose dir should be pruned: %v", err)
	}
	if _, err := os.Stat(filepath.Join(genRoot, KindBMaterializedList)); !os.IsNotExist(err) {
		t.Fatalf("materialized list should be removed")
	}
}

func TestCleanupKindBMaterialized_deepNestedPrune(t *testing.T) {
	t.Parallel()
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, DoctestInternalExposeDir, "a", "b", "c", "d", "expose.go")
	if err := os.MkdirAll(filepath.Dir(virt), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(virt, []byte("package d\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := recordKindBMaterialized(genRoot, virt); err != nil {
		t.Fatal(err)
	}
	if err := CleanupKindBMaterialized(genRoot); err != nil {
		t.Fatalf("cleanup: %v", err)
	}
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("expose.go still present: %v", err)
	}
	if _, err := os.Stat(filepath.Join(product, DoctestInternalExposeDir)); !os.IsNotExist(err) {
		t.Fatalf("deep expose dir should be pruned: %v", err)
	}
}

func TestCleanupKindBMaterialized_refusesNonExposePath(t *testing.T) {
	t.Parallel()
	genRoot := t.TempDir()
	secret := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(secret, []byte("keep\n"), 0644); err != nil {
		t.Fatal(err)
	}
	listPath := filepath.Join(genRoot, KindBMaterializedList)
	if err := os.WriteFile(listPath, []byte(secret+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := CleanupKindBMaterialized(genRoot); err == nil {
		t.Fatal("expected error refusing non-expose path")
	}
	if _, err := os.Stat(secret); err != nil {
		t.Fatalf("non-expose file must not be removed: %v", err)
	}
	data, err := os.ReadFile(listPath)
	if err != nil {
		t.Fatalf("list should be kept for leftover paths: %v", err)
	}
	if !strings.Contains(string(data), secret) {
		t.Fatalf("list should still name leftover path; got %q", data)
	}
}

func TestCleanupKindBMaterialized_keepsListOnRemoveFailure(t *testing.T) {
	t.Parallel()
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(virt, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(virt, "keep.txt"), []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := recordKindBMaterialized(genRoot, virt); err != nil {
		t.Fatal(err)
	}
	if err := CleanupKindBMaterialized(genRoot); err == nil {
		t.Fatal("expected error when expose.go is a non-empty dir")
	}
	if _, err := os.Stat(virt); err != nil {
		t.Fatalf("non-empty expose.go dir should remain: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(genRoot, KindBMaterializedList))
	if err != nil {
		t.Fatalf("list should be kept after remove failure: %v", err)
	}
	if !strings.Contains(string(data), virt) {
		t.Fatalf("list should still name leftover path; got %q", data)
	}
}

func TestRecordKindBMaterialized_rejectsNonExpose(t *testing.T) {
	t.Parallel()
	genRoot := t.TempDir()
	if err := recordKindBMaterialized(genRoot, filepath.Join(t.TempDir(), "other.go")); err == nil {
		t.Fatal("expected reject of non-expose path")
	}
	if _, err := os.Stat(filepath.Join(genRoot, KindBMaterializedList)); !os.IsNotExist(err) {
		t.Fatalf("must not record rejected path")
	}
}

// hasVarReexport reports whether body re-exports name as a var alias of the
// imported internal package (any import alias).
func hasVarReexport(body, name string) bool {
	// e.g. var DefaultName = srcpkg.DefaultName
	//      var DefaultName = src.DefaultName
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "var "+name) {
			continue
		}
		if strings.Contains(line, "."+name) {
			return true
		}
	}
	return false
}

// hasConstOrVarReexport accepts either const MaxN = alias.MaxN or var MaxN = alias.MaxN.
func hasConstOrVarReexport(body, name string) bool {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "const "+name) || strings.HasPrefix(line, "var "+name) {
			if strings.Contains(line, "."+name) {
				return true
			}
		}
	}
	return false
}
