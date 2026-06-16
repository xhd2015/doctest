# Scenario

**Feature**: the format string contains `%s` verbs with corresponding arguments, and already ends with `\n`

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- The format string contains `%s` verbs with corresponding arguments, and already ends with `\n`.

## Steps
1. Set `LOGF_FORMAT=item=%s count=%s\n` via env.
2. Set `req.Args = ["alpha", "42"]` as format arguments.
3. Call `Logf("item=%s count=%s\n", "alpha", "42")` via the parent Run.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "LOGF_FORMAT=item=%s count=%s\n")
    req.Args = []string{"alpha", "42"}
    return nil
}
```
