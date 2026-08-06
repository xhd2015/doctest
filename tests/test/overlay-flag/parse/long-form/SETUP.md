# Scenario

**Feature**: long form `--overlay FILE` is accepted and stored on Options

```
ParseTestOptions(["--overlay", "/abs/user.json", "."])
  -> err=nil, Options.Overlay == /abs/user.json
```

## Preconditions

- Absolute path; does not re-test relative resolution (see `relative-abs`).
- Classic TDD: RED until `--overlay` is a known parse flag.

## Steps

1. Write an absolute overlay path under session temp.
2. Parse with `--overlay`.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(os.TempDir(), "doctest-overlay-flag-"+d.DOCTEST_SESSION_ID, "long")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	abs := filepath.Join(dir, "user-overlay.json")
	if err := os.WriteFile(abs, []byte(`{"Replace":{}}`), 0o644); err != nil {
		t.Fatalf("write overlay: %v", err)
	}
	req.ParseArgs = []string{"--overlay", abs, "."}
	return nil
}
```
