# Scenario

**Feature**: parent of cwd is not displayed as `./..`; falls through to home or absolute

```
# cwd rules fail for parent
strict child of cwd -> "./" + rel   # parent is NOT a strict child

# home shorten or fallback
path under home -> "~" + suffix | otherwise -> absolute unchanged
```

## Preconditions

- Cwd is a nested subdirectory; target path is the parent directory.

## Steps

1. Create a nested cwd under home (not under a temp dir outside home).
2. `chdir` to the nested directory.
3. Set `req.Path` to the parent directory.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	base := filepath.Join(home, ".doctest-display-parent-test")
	mkdirAll(t, base)
	nested := filepath.Join(base, "nested")
	mkdirAll(t, nested)
	t.Cleanup(func() { _ = os.RemoveAll(base) })
	chdirTo(t, nested)
	req.Path = base
	return nil
}
```