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

func TestDoctestInternalExposeDir(t *testing.T) {
	t.Parallel()
	if DoctestInternalExposeDir != "__doctest_internal_expose" {
		t.Fatalf("DoctestInternalExposeDir = %q", DoctestInternalExposeDir)
	}
	if DoctestInternalShimDir != "__doctest_internal_shim" {
		t.Fatalf("DoctestInternalShimDir = %q", DoctestInternalShimDir)
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
