package build

import (
	"reflect"
	"testing"
)

func TestPackageTestShardsDistributesAllPackages(t *testing.T) {
	pkgs := []string{"a", "b", "c", "d", "e"}
	shards := packageTestShards(pkgs, 3)
	if len(shards) != 3 {
		t.Fatalf("shards=%d want 3: %v", len(shards), shards)
	}
	seen := map[string]bool{}
	for _, s := range shards {
		for _, p := range s {
			if seen[p] {
				t.Fatalf("duplicate package %q", p)
			}
			seen[p] = true
		}
	}
	if len(seen) != len(pkgs) {
		t.Fatalf("seen %d packages, want %d", len(seen), len(pkgs))
	}
}

func TestPackageTestShardsSerialWhenWorkersOne(t *testing.T) {
	pkgs := []string{"a", "b"}
	got := packageTestShards(pkgs, 1)
	want := [][]string{pkgs}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}
