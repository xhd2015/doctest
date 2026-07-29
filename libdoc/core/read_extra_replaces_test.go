package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadExtraReplacesModuleAndPath(t *testing.T) {
	dir := t.TempDir()
	mod := `module example.com/app

go 1.21

replace github.com/micro/go-micro => github.com/sixhuan/go-micro v1.10.1
replace github.com/gansidui/go-utils => ./tools/work_around/github.com/gansidui/go-utils
replace (
	golang.org/x/net => golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f
)
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(mod), 0644); err != nil {
		t.Fatal(err)
	}
	out, parentPath, parentModule := readExtraReplaces(dir, "example.com/app")
	if !strings.Contains(out, "replace github.com/micro/go-micro => github.com/sixhuan/go-micro v1.10.1") {
		t.Fatalf("missing module→module replace:\n%s", out)
	}
	if !strings.Contains(out, "replace github.com/gansidui/go-utils => ") {
		t.Fatalf("missing path replace:\n%s", out)
	}
	if !strings.Contains(out, filepath.Join(dir, "tools/work_around/github.com/gansidui/go-utils")) {
		t.Fatalf("path not absolutized:\n%s", out)
	}
	if !strings.Contains(out, "replace golang.org/x/net => golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f") {
		t.Fatalf("missing block replace:\n%s", out)
	}
	// Path vs module→module are tracked separately for vendor inject policy.
	if !parentPath["github.com/gansidui/go-utils"] {
		t.Fatalf("expected path-replaced module gansidui in path set: %v", parentPath)
	}
	if parentModule["github.com/gansidui/go-utils"] {
		t.Fatalf("path replace must not be in module set: %v", parentModule)
	}
	if !parentModule["github.com/micro/go-micro"] {
		t.Fatalf("expected module→module replace micro in module set: %v", parentModule)
	}
	if parentPath["github.com/micro/go-micro"] {
		t.Fatalf("module→module must not be in path set: %v", parentPath)
	}
	if !parentModule["golang.org/x/net"] {
		t.Fatalf("expected block module→module replace x/net in module set: %v", parentModule)
	}
}
