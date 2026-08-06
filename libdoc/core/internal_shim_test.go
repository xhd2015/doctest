package core

import "testing"

func TestKindAShimImport(t *testing.T) {
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
	mod, full, ok := productInternalImport("example.com/app/internal/greet", "example.com/runner")
	if !ok || mod != "example.com/app" || full != "example.com/app/internal/greet" {
		t.Fatalf("got mod=%q full=%q ok=%v", mod, full, ok)
	}
	// Parent-module internal is internal-compile territory.
	if _, _, ok := productInternalImport("example.com/app/internal/greet", "example.com/app"); ok {
		t.Fatal("parent internal should not be treated as kind B expose")
	}
	if _, _, ok := productInternalImport("testcase/x/internal/y", ""); ok {
		t.Fatal("testcase internal is kind A not B")
	}
}
