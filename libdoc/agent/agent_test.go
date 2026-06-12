package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCreatesExpectedFiles(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "generated")
	err := Generate(GenerateOptions{
		Idea: "a multi word idea",
		Dir:  outDir,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	for _, name := range []string{"DOCTEST.md", "SETUP.md"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("%s missing: %v", name, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "SETUP.md"))
	if err != nil {
		t.Fatalf("read SETUP.md: %v", err)
	}
	if !strings.Contains(string(data), "Generated from idea: a multi word idea") {
		t.Fatalf("SETUP.md missing idea:\n%s", string(data))
	}
}

func TestGenerateDefaultDirUsesWorkingDirectory(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	if err := Generate(GenerateOptions{Idea: "default dir"}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmp, "doctest-test-cases", "DOCTEST.md")); err != nil {
		t.Fatalf("default README missing: %v", err)
	}
}

func TestGenerateCreatesParentDirectories(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "missing", "parents", "generated")
	if err := Generate(GenerateOptions{Idea: "parents", Dir: outDir}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "SETUP.md")); err != nil {
		t.Fatalf("generated setup missing: %v", err)
	}
}

func TestFillCode(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(t.TempDir(), "not-dir")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		dir     string
		wantErr string
	}{
		{name: "existing dir", dir: dir},
		{name: "blank dir", dir: "   ", wantErr: "agent fill-code requires <target-dir>"},
		{name: "missing dir", dir: filepath.Join(t.TempDir(), "missing"), wantErr: "no such file"},
		{name: "file instead of dir", dir: file, wantErr: "not a directory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FillCode(tt.dir)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("FillCode(%q): %v", tt.dir, err)
				}
				return
			}
			if err == nil {
				t.Fatalf("FillCode(%q): expected error containing %q", tt.dir, tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("FillCode(%q): expected error containing %q, got %v", tt.dir, tt.wantErr, err)
			}
		})
	}
}
