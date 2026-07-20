# Scenario

**Feature**: top-level `session.Once` caches JSON once per (session id, key)

```
# caller invokes Once with session env and key
Caller -> session.Once(t, key, fn)
session.Once -> syscall.Getenv(DOCTEST_SESSION_ID)
session.Once -> t.TempDir()/session-once/<slug>/  (not UserCacheDir; keeps go testcache warm)

# success path
fn(cacheDir) -> json.RawMessage
session.Once -> processMemo + optional temp value file under lock
Caller <- same json.RawMessage on later Once same key (in-process)

# error path
fn -> error
session.Once -> processMemo error; later Once returns error without re-running fn
```

## Preconditions

- Public package path is `github.com/xhd2015/doctest/session` (top-level, like `assert`).
- API under test:
  `func Once(t testing.TB, key string, fn func(t testing.TB, cacheDir string) (json.RawMessage, error)) (json.RawMessage, error)`.
- Session id is read only via `syscall.Getenv("DOCTEST_SESSION_ID")`.
- On-disk scratch (optional):
  `t.TempDir()/session-once/<slug(key)>/{lock,value,error}` — under test temp so
  opens are not go testcache inputs.
- Cross-call reuse in one process is primarily `processMemo` (session id + key).
- `value` stores **raw JSON bytes** (not a plain string path).

## Steps

1. Leaf `Setup` sets `req.SessionID`, `req.Key`, `req.Mode`, and call flags.
2. Root `Run` applies session env with `t.Setenv`, invokes `session.Once`, and
   records values, errors, fn call count, and optional cache file bytes.
3. Leaf `Assert` checks errors, JSON equality, layout, or unmarshal as specified.

## Context

- Each success leaf uses a unique `SessionID` so parallel leaves do not share
  once-dirs.
- Harness must not call `os.Getenv("DOCTEST_SESSION_ID")` for product behavior;
  product uses `syscall.Getenv` only.
- `DOCTEST_SESSION_ID` injected variable is for harness session-scoped caching
  only; product session id is controlled via `req.SessionID` + `t.Setenv`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Defaults; leaves override.
	if req.Mode == "" {
		req.Mode = "json-object"
	}
	return nil
}

// sessionOnceSegment is the path segment under t.TempDir() used by Once.
const sessionOnceSegment = "session-once"
```
