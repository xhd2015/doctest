package core

import (
	"strings"
	"testing"
)

func TestReconcileGeneratedImportsDoesNotAutoAddStdlibAndDropsUnused(t *testing.T) {
	src := `package pkg

import (
	droot "testcase/__droot"
	"testing"
)

func Setup(t *testing.T, req *droot.Request) error {
	_ = 20 * time.Second
	_ = req
	return nil
}
`
	out, err := reconcileGeneratedImports(src)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, `"time"`) {
		t.Fatalf("must not auto-add time import from time. selector:\n%s", out)
	}
	// droot is used via droot.Request
	if !strings.Contains(out, "droot") {
		t.Fatalf("expected droot kept:\n%s", out)
	}
	// time. selector preserved even without import (compile will fail later)
	if !strings.Contains(out, "time.Second") {
		t.Fatalf("expected time.Second body preserved:\n%s", out)
	}

	src2 := `package pkg

import (
	parent "testcase/parent"
	"testing"
)

func Setup(t *testing.T) error {
	_ = t
	return nil
}
`
	out2, err := reconcileGeneratedImports(src2)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out2, "parent") {
		t.Fatalf("expected unused parent import removed:\n%s", out2)
	}
	if !strings.Contains(out2, `"testing"`) {
		t.Fatalf("expected testing kept:\n%s", out2)
	}
}

func TestFormatGeneratedGoNoFormatSourceOnParseFail(t *testing.T) {
	// Invalid Go: reconcile fails; formatGeneratedGo returns bytes as-is, no error.
	src := []byte("package pkg\n\nfunc broken( {\n")
	out, err := formatGeneratedGo("x.go", src)
	if err != nil {
		t.Fatalf("formatGeneratedGo must not error on parse fail: %v", err)
	}
	if string(out) != string(src) {
		t.Fatalf("expected original bytes on parse fail, got:\n%s", out)
	}
}
