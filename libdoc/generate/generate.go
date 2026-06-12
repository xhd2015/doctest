package generate

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	opencoderun "github.com/xhd2015/agent-pro/agent/opencode/run"
)

//go:embed PROMPT.md
var SystemPrompt string

type GenerateOptions struct {
	DryRun bool
	Model  string
}

type FileNeed struct {
	Path    string
	IsSetup bool
	IsRoot  bool
	Content string
}

func Scan(root string) ([]FileNeed, error) {
	var needs []FileNeed
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		name := d.Name()
		if name != "SETUP.md" && name != "ASSERT.md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		if !hasGoBlock(content) {
			rel, _ := filepath.Rel(root, path)
			isRoot := name == "SETUP.md" && rel == "SETUP.md"
			needs = append(needs, FileNeed{
				Path:    rel,
				IsSetup: name == "SETUP.md",
				IsRoot:  isRoot,
				Content: content,
			})
		}
		return nil
	})
	return needs, err
}

func BuildPrompt(need FileNeed) string {
	var b strings.Builder
	b.WriteString(SystemPrompt)
	b.WriteString(fmt.Sprintf("\n\n## File: %s\n", need.Path))
	if need.IsRoot {
		b.WriteString("This is the ROOT SETUP.md.\n")
	}
	b.WriteString(fmt.Sprintf("\n## Current content:\n\n```markdown\n%s\n```\n", need.Content))
	b.WriteString("\nFill in the executable Go code block. Output the COMPLETE new file content.")
	return b.String()
}

func Generate(prompt string, model string) (string, error) {
	output, _, err := opencoderun.Run(context.Background(), opencoderun.Options{
		Prompt: prompt,
		Model:  model,
		Dir:    ".",
	})
	return output, err
}

func Apply(root string, need FileNeed, result string) error {
	result = strings.TrimSpace(result)
	path := filepath.Join(root, need.Path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(result+"\n"), 0644)
}

func Run(root string, opts GenerateOptions) error {
	needs, err := Scan(root)
	if err != nil {
		return err
	}
	if len(needs) == 0 {
		fmt.Println("All files already have Go code blocks.")
		return nil
	}
	for _, need := range needs {
		if opts.DryRun {
			fmt.Printf("Would generate Go code for %s\n", need.Path)
			continue
		}
		result, err := Generate(BuildPrompt(need), opts.Model)
		if err != nil {
			return fmt.Errorf("%s: %w", need.Path, err)
		}
		if err := Apply(root, need, result); err != nil {
			return err
		}
		fmt.Printf("Generated Go code for %s\n", need.Path)
	}
	return nil
}

func hasGoBlock(content string) bool {
	i := 0
	for {
		start := strings.Index(content[i:], "```go")
		if start < 0 {
			return false
		}
		start += i
		lineEnd := strings.IndexByte(content[start:], '\n')
		if lineEnd < 0 {
			return false
		}
		codeStart := start + lineEnd + 1
		close := strings.Index(content[codeStart:], "```")
		if close < 0 {
			return false
		}
		close += codeStart
		end := close + len("```")
		if end < len(content) && content[end] == '\r' {
			end++
		}
		if end < len(content) && content[end] == '\n' {
			end++
		}
		if strings.TrimSpace(content[close+len("```"):]) == "" {
			return true
		}
		i = end
	}
}
