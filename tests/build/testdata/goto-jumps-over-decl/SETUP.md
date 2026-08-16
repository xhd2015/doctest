# Scenario

**Feature**: SETUP helper uses `goto mock`; build-path `:=` is in an inner block

```
# helper copies a binary then goto mock; else { cmd := exec.Command }
# generated suite must compile that helper
```

## Setup

- Copy-success path still `goto mock`.
- Build path `:=` lives in an inner block so `goto` does not skip a declaration.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
)

func installFakeOpencode(bin string) error {
	if p, err := exec.LookPath("true"); err == nil && p != "" {
		if data, rerr := os.ReadFile(p); rerr == nil {
			if werr := os.WriteFile(bin, data, 0o755); werr == nil {
				goto mock
			}
		}
	}
	{
		cmd := exec.Command("true")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		_ = out
	}
mock:
	_ = filepath.Join(".", "fake-opencode.json")
	return nil
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Name = "goto-jumps-over-decl"
	return nil
}
```
