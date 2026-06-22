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
- This is a standalone root with its own `DOCTEST.md` because `subagent.Logf` is called in-process.

## Steps
1. Read `LOGF_FORMAT` from `req.Env` as the format string (default: `"default"`).
2. Read `req.Args` as variadic format arguments (all strings).
3. Redirect `os.Stdout` to a pipe, call `subagent.Logf(format, args...)`.
4. Restore `os.Stdout`, read captured output, return as `resp.Stdout`.

```go
import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"github.com/xhd2015/agent-pro/agent/subagent"
)

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "TEST_GROUP=logf")
	return nil
}
```
