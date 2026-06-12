package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanFindsFilesWithoutGoBlocks(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "SETUP.md", "# Setup\n\nNeed a Go block.\n")
	writeFile(t, root, "leaf/SETUP.md", "# Setup\n\nAlso needs a Go block.\n")
	writeFile(t, root, "leaf/ASSERT.md", "# Assert\n\nNeeds a Go block too.\n")
	writeFile(t, root, "leaf/README.md", "# Not a SETUP or ASSERT\n")

	needs, err := Scan(root)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(needs) != 3 {
		t.Fatalf("expected 3 needs, got %d: %+v", len(needs), needs)
	}

	found := make(map[string]bool)
	for _, n := range needs {
		found[n.Path] = true
	}
	if !found["SETUP.md"] || !found["leaf/SETUP.md"] || !found["leaf/ASSERT.md"] {
		t.Fatalf("missing expected paths, got %v", needs)
	}
}

func TestScanSkipsFilesWithGoBlocks(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "SETUP.md", "# Setup\n\n```go\ntype Request struct{}\n```\n")
	writeFile(t, root, "leaf/ASSERT.md", "# Assert\n\n```go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n```\n")

	needs, err := Scan(root)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(needs) != 0 {
		t.Fatalf("expected 0 needs for files with Go blocks, got %d", len(needs))
	}
}

func TestScanSkipsTestdataDir(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "SETUP.md", "# Setup\n\nNeed a Go block.\n")
	writeFile(t, root, "testdata/SETUP.md", "# Setup\n\nNeed a Go block.\n")

	needs, err := Scan(root)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(needs) != 1 {
		t.Fatalf("expected 1 need (testdata skipped), got %d", len(needs))
	}
	if len(needs) > 0 && needs[0].Path != "SETUP.md" {
		t.Fatalf("expected root SETUP.md, got %q", needs[0].Path)
	}
}

func TestScanIdentifiesRootSetup(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "SETUP.md", "# Setup\n\nNeeds a Go block.\n")
	writeFile(t, root, "leaf/SETUP.md", "# Setup\n\nAlso needs a Go block.\n")

	needs, err := Scan(root)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	for _, n := range needs {
		if n.Path == "SETUP.md" && !n.IsRoot {
			t.Fatalf("expected root SETUP.md to be IsRoot")
		}
		if n.Path == "leaf/SETUP.md" && n.IsRoot {
			t.Fatalf("expected leaf SETUP.md to NOT be IsRoot")
		}
	}
}

func TestBuildPromptIncludesFileContent(t *testing.T) {
	need := FileNeed{
		Path:    "SETUP.md",
		IsSetup: true,
		IsRoot:  true,
		Content: "# Setup\n\nProse only.\n",
	}
	prompt := BuildPrompt(need)
	if !strings.Contains(prompt, "This is the ROOT SETUP.md") {
		t.Fatalf("expected root annotation, got %q", prompt)
	}
	if !strings.Contains(prompt, "# Setup\n\nProse only.") {
		t.Fatalf("expected file content in prompt, got %q", prompt)
	}
	if !strings.Contains(prompt, SystemPrompt) {
		t.Fatalf("expected system prompt in output, got %q", prompt)
	}
}

func TestBuildPromptNonRoot(t *testing.T) {
	need := FileNeed{
		Path:    "leaf/ASSERT.md",
		IsSetup: false,
		IsRoot:  false,
		Content: "# Assert\n\nProse only.\n",
	}
	prompt := BuildPrompt(need)
	if strings.Contains(prompt, "This is the ROOT SETUP.md") {
		t.Fatal("expected no root annotation for non-root file")
	}
	if !strings.Contains(prompt, "File: leaf/ASSERT.md") {
		t.Fatalf("expected file path in prompt, got %q", prompt)
	}
}

func TestRunDryRunDoesNotWrite(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "SETUP.md", "# Setup\n\nNeeds a Go block.\n")

	if err := Run(root, GenerateOptions{DryRun: true}); err != nil {
		t.Fatalf("Run dry-run: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "SETUP.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "type Request") {
		t.Fatalf("dry run should not write, got:\n%s", data)
	}
}

func writeFile(t *testing.T, root string, rel string, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
