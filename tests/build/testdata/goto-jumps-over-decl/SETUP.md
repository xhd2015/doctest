# Scenario

**Feature**: SETUP helper uses goto over a later `:=` (wrk unwind-pipeline / fake-opencode)

```
# helper copies a binary then goto mock; else declares cmd := exec.Command
# generated suite must compile that helper
```

## Setup

- Mirrors `installFakeOpencode`: early `goto mock` skips `cmd := exec.Command`.
- Go rejects `goto` over a new variable (`jumps over declaration of cmd`).

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
	cmd := exec.Command("true")
	if out, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	_ = out
mock:
	_ = filepath.Join(".", "fake-opencode.json")
	return nil
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Name = "goto-jumps-over-decl"
	return nil
}
```
