# Scenario

**Feature**: the target path exists but is not a directory (**L2 short path**)

```
doctest agent fill-code <file-not-dir>
  -> non-zero; stderr mentions "not a directory"
```

## Preconditions

- Path validation fails before any agent/LLM work — in-process CLI.

## Steps

1. Create a file path (not a directory).
2. Run `doctest agent fill-code <file>` in-process.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	path := filepath.Join(t.TempDir(), "target.txt")
	if err := os.WriteFile(path, []byte("not a dir"), 0644); err != nil {
		t.Fatal(err)
	}
	// L2 short path: no Env isolation, no product binary.
	req.UseCLI = false
	req.Bin = ""
	req.Env = nil
	req.Args = []string{"agent", "fill-code", path}
	return nil
}
```
