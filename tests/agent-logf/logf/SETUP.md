# Scenario

**Feature**: the `Logf` function is in `agent-pro/agent/subagent` and writes timestamped output to `os.Stdout`

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- The `Logf` function is in `agent-pro/agent/subagent` and writes timestamped output to `os.Stdout`.
- Each leaf provides a format string and optional args via environment variables.
- This is a standalone root with its own `DOCTEST.md` because `subagent.Logf` is probed via subprocess.
- Parent never reassigns `os.Stdout`; capture is child-only.

## Steps
1. Resolve module root onto `req.ModuleRoot` for `go list` / probe module replace.
2. Read `LOGF_FORMAT` from `req.Env` as the format string (default: `"default"`).
3. Read `req.Args` as variadic format arguments (all strings).
4. Run a temp `go run` probe that calls `subagent.Logf`; capture child stdout.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Env = append(req.Env, "TEST_GROUP=logf")
	// logf/ -> agent-logf/ -> tests/ -> module root
	req.ModuleRoot = filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
	return nil
}
```
