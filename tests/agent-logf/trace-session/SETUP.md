## Preconditions
- A doctest binary is available at `req.Bin` (built by the root SETUP.md).
- Session directories are created under a temp `DOCTEST_DEBUG_SESSION_HOME`.

## Steps
1. Each leaf creates a session directory with `meta.json` and events.
2. The root `Run` shells out to `doctest agent implement --trace --session-id X`.
3. Leaves assert the output format.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=trace-session")
    return nil
}
```
