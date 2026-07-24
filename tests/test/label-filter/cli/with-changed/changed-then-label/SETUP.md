# Scenario

**Feature**: changed-file filter applies before label filter

```
change slow only; --changed --label 'slow || heavy' -> runs slow only
```

## Steps

1. Commit baseline mod; unstaged edit on `slow/ASSERT.md`.
2. Run from repo root with both flags.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	repoDir, modDir := createLabelFilterGitMod(t)
	assertPath := filepath.Join(modDir, "slow", "ASSERT.md")
	data, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, []byte("\n<!-- changed -->\n")...)
	if err := os.WriteFile(assertPath, data, 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = repoDir
	req.Args = []string{"test", modDir, "--changed", "--label", "slow || heavy"}
	return nil
}
```