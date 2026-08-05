package build

import (
	"reflect"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func TestAppendOptsGoTestFlagsCoverAndExec(t *testing.T) {
	par := 2
	opts := core.Options{
		Cover:     true,
		CoverMode: "atomic",
		CoverPkg:  "example.com/mod/...",
		Race:      true,
		Short:     true,
		FailFast:  true,
		Parallel:  &par,
		Shuffle:   "on",
		Tags:      "integration",
		Gcflags:   "all=-N",
		Ldflags:   "-X main.v=1",
	}
	got := appendOptsGoTestFlags(nil, opts)
	wantSub := []string{
		"-covermode=atomic",
		"-coverpkg=example.com/mod/...",
		"-cover",
		"-race",
		"-short",
		"-failfast",
		"-parallel=2",
		"-shuffle=on",
		"-tags=integration",
		"-gcflags=all=-N",
		"-ldflags=-X main.v=1",
	}
	for _, w := range wantSub {
		if !containsArg(got, w) {
			t.Errorf("missing %q in %v", w, got)
		}
	}
}

func TestAppendOptsGoTestFlagsCoverModeImpliesCover(t *testing.T) {
	got := appendOptsGoTestFlags(nil, core.Options{CoverMode: "set"})
	if !containsArg(got, "-covermode=set") || !containsArg(got, "-cover") {
		t.Fatalf("got %v", got)
	}
}

func TestCheckCoverProfilePackages(t *testing.T) {
	if err := checkCoverProfilePackages(core.Options{CoverProfile: "/tmp/c.out"}, []string{"./suite"}); err != nil {
		t.Fatal(err)
	}
	err := checkCoverProfilePackages(core.Options{CoverProfile: "/tmp/c.out"}, []string{"./a", "./b"})
	if err == nil {
		t.Fatal("expected multi-package error")
	}
	if !strings.Contains(err.Error(), "-coverprofile") || !strings.Contains(err.Error(), "multiple packages") {
		t.Fatalf("err=%v", err)
	}
	if err := checkCoverProfilePackages(core.Options{}, []string{"./a", "./b"}); err != nil {
		t.Fatal(err)
	}
}

func TestLeafCacheMeasureNoSkip(t *testing.T) {
	if (core.Options{}).LeafCacheMeasureNoSkip() {
		t.Fatal("default should allow skip")
	}
	cases := []core.Options{
		{Cover: true},
		{CoverProfile: "x"},
		{CoverMode: "set"},
		{CoverPkg: "p"},
		{Race: true},
	}
	for _, o := range cases {
		if !o.LeafCacheMeasureNoSkip() {
			t.Fatalf("expected measure no-skip for %#v", o)
		}
	}
}

func containsArg(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}

func TestAppendOptsGoTestFlagsEmpty(t *testing.T) {
	got := appendOptsGoTestFlags([]string{"test"}, core.Options{})
	want := []string{"test"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
