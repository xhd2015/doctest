package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunValidationCases(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr string
	}{
		{
			name: "valid root only",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "setup\n")
				return dir
			},
		},
		{
			name: "valid leaf",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "setup\n")
				writeFile(t, dir, "leaf/SETUP.md", "setup\n")
				writeFile(t, dir, "leaf/ASSERT.md", "assert\n")
				return dir
			},
		},
		{
			name: "missing readme",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "SETUP.md", "setup\n")
				return dir
			},
			wantErr: "root must contain DOCTEST.md",
		},
		{
			name: "file instead of dir",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				path := filepath.Join(dir, "file.md")
				writeFile(t, dir, "file.md", "file\n")
				return path
			},
			wantErr: "not a directory",
		},
		{
			name: "assert without setup at root",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "ASSERT.md", "assert\n")
				return dir
			},
			wantErr: "ASSERT.md found but SETUP.md missing",
		},
		{
			name: "assert without setup in leaf",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "setup\n")
				writeFile(t, dir, "leaf/ASSERT.md", "assert\n")
				return dir
			},
			wantErr: "ASSERT.md found but SETUP.md missing",
		},
		{
			name: "testdata skipped",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "setup\n")
				writeFile(t, dir, "testdata/bad/ASSERT.md", "assert\n")
				return dir
			},
		},
		{
			name: "extra files allowed",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "setup\n")
				writeFile(t, dir, "notes.txt", "notes\n")
				writeFile(t, dir, "leaf/SETUP.md", "setup\n")
				writeFile(t, dir, "leaf/ASSERT.md", "assert\n")
				writeFile(t, dir, "leaf/notes.txt", "notes\n")
				return dir
			},
		},
		{
			name: "embedded go program in raw string",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.MainGo = `package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}`\n\treturn nil\n}\n```\n")
				return dir
			},
			wantErr: "anti-pattern: raw Go code embedded in string literal",
		},
		{
			name: "embedded go program in interpreted string",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.MainGo = \"package main\\n\\nfunc main() {}\\n\"\n\treturn nil\n}\n```\n")
				return dir
			},
			wantErr: "anti-pattern: raw Go code embedded in string literal",
		},
		{
			name: "go test shell-out",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "# Setup\n\n```go\nimport \"os/exec\"\n\nfunc Run(t *testing.T, req *Request) (*Response, error) {\n\tcmd := exec.Command(\"go\", \"test\", \"./pkg/foo\", \"-run\", \"TestFoo\")\n\tout, _ := cmd.CombinedOutput()\n\treturn &Response{Stdout: string(out)}, nil\n}\n```\n")
				return dir
			},
			wantErr: "anti-pattern: shelling out to 'go test'",
		},
		{
			name: "valid String without both keywords",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.Msg = \"package main is the entry point\"\n\treturn nil\n}\n```\n")
				return dir
			},
		},
		{
			name: "valid exec.Command not go test",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "# Setup\n\n```go\nimport \"os/exec\"\n\nfunc Run(t *testing.T, req *Request) (*Response, error) {\n\tcmd := exec.Command(\"doctest\", \"build\", req.InputDir)\n\tout, _ := cmd.CombinedOutput()\n\treturn &Response{Stdout: string(out)}, nil\n}\n```\n")
				return dir
			},
		},
		{
			name: "anti-pattern in testdata skipped",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "## Setup\n")
				writeFile(t, dir, "testdata/bad/SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.MainGo = `package main\nfunc main() {}`\n\treturn nil\n}\n```\n")
				return dir
			},
		},
		{
			name: "multiple anti-patterns in one tree",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				writeFile(t, dir, "DOCTEST.md", "# tests\n")
				writeFile(t, dir, "SETUP.md", "# Setup\n\n```go\nfunc Run(t *testing.T, req *Request) (*Response, error) {\n\tcmd := exec.Command(\"go\", \"test\", \"./pkg/foo\")\n\tout, _ := cmd.CombinedOutput()\n\treturn &Response{Stdout: string(out)}, nil\n}\n```\n")
				writeFile(t, dir, "leaf/SETUP.md", "# Setup\n\n```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.MainGo = `package main\nfunc main() {}`\n\treturn nil\n}\n```\n")
				writeFile(t, dir, "leaf/ASSERT.md", "# Assert\n\n```go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n```\n")
				return dir
			},
			wantErr: "anti-pattern: shelling out to 'go test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Run(tt.setup(t))
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Run: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
