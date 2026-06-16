# Scenario

**Feature**: a fresh sessions directory with no session subdirectories

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A fresh sessions directory with no session subdirectories.
- The `--list-sessions` flag should report no sessions found.

## Steps
1. Ensure the sessions home directory exists but is empty.
2. Run `doctest agent implement --list-sessions`.

```go
import (
    "os"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    sessHome := sessionsDir()
    os.MkdirAll(sessHome, 0755)
    os.RemoveAll(sessHome)
    os.MkdirAll(sessHome, 0755)

    req.Args = []string{"agent", "implement", "--list-sessions"}
    return nil
}
```
