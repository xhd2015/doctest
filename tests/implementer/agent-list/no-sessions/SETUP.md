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
