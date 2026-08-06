# Scenario

**Feature**: relative `-overlay` path is abs-resolved against process cwd at parse

```
ParseTestOptions(["-overlay", "rel/user-overlay.json", "."])
  -> Options.Overlay is absolute and equals Abs(rel)
  -> relative form must not remain unresolved
```

## Preconditions

- Relative path only (no leading `/`).
- Resolution class matches profile flags (`filepath.Abs` / process cwd).
- File need not exist for parse (same as many path flags); we still create a
  real relative file under a known cwd-adjacent temp when practical.

## Steps

1. Create a relative path string (not absolute).
2. Parse with `-overlay` + relative path.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Use a relative path that is stable: write under TempDir then pass a path
	// relative to process cwd is hard without Chdir. Instead pass a clearly
	// relative token "overlay-flag-rel-<session>/user.json" and create it under
	// cwd via Abs for the file write only — parse still receives the relative form.
	rel := filepath.Join("overlay-flag-rel-"+d.DOCTEST_SESSION_ID, "user.json")
	if filepath.IsAbs(rel) {
		t.Fatalf("fixture path must be relative, got %q", rel)
	}
	abs, err := filepath.Abs(rel)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(abs, []byte(`{"Replace":{}}`), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Dir(abs))
	})
	req.ParseArgs = []string{"-overlay", rel, "."}
	return nil
}
```
