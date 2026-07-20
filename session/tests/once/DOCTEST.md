# Session Once — json.RawMessage API

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Caller** — integration leaf or helper (e.g. `testbin.Ensure`) that needs a
  value computed at most once per `doctest test` run on this machine.
- **Session Once** (`github.com/xhd2015/doctest/session`) — public package API
  `Once(t, key, fn) (json.RawMessage, error)`.
- **Session environment** — process env `DOCTEST_SESSION_ID`, read only via
  `syscall.Getenv` (never `os.Getenv`, so Go's test result cache is not keyed
  on the session id).
- **Disk scratch** — under `t.TempDir()/session-once/<slug(key)>/` (not
  UserCacheDir) so opens do not bust go package testcache when session id
  changes; files `lock`, `value`, `error` as needed.
- **processMemo** — in-process reuse of success/error for `(session, key)`.
- **Fn** — user callback `func(t testing.TB, cacheDir string) (json.RawMessage, error)`
  that receives a writable `cacheDir` and returns JSON bytes on success.

**Behaviors**

- Missing or empty `DOCTEST_SESSION_ID` → error; `fn` is not run.
- Empty `key` → error; `fn` is not run.
- First success for `(session, key)` runs `fn` once, memos the JSON, returns
  the same bytes to every later caller with the same pair **in this process**.
- Different keys under the same session are independent (separate `fn` invocations).
- If `fn` returns an error, that error is memoized; a second `Once` with the
  same key returns the error **without** re-running `fn`.
- Returned `json.RawMessage` is valid JSON the client can `json.Unmarshal` into
  a struct (e.g. `{"path":"..."}` for binary paths).

## Decision Tree

```
once/
├── validation/                         [inputs rejected before fn]
│   ├── missing-session-id/             V1: no DOCTEST_SESSION_ID
│   └── empty-key/                      V2: key == ""
├── success/                            [sid + key valid; fn succeeds]
│   ├── first-write-and-reuse/          S1: write JSON; second call caches
│   ├── different-keys-independent/     S2: key A and B independent
│   ├── cache-dir-layout/               S3: path shape + writable
│   └── json-unmarshal-struct/          S4: Unmarshal into struct
└── error/                              [fn fails]
    └── error-persisted/                E1: error replayed; fn once
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `validation/missing-session-id` | V1 — empty/missing session id → error, fn not called |
| `validation/empty-key` | V2 — empty key → error, fn not called |
| `success/first-write-and-reuse` | S1 — first writes JSON object; second same key: fn once, raw equal |
| `success/different-keys-independent` | S2 — different keys each invoke fn once with distinct values |
| `success/cache-dir-layout` | S3 — cacheDir is under `session-once/...` (temp) and writable |
| `success/json-unmarshal-struct` | S4 — returned bytes `json.Unmarshal` into a typed struct |
| `error/error-persisted` | E1 — fn error persisted; second Once returns error without success |

## How to Run

```sh
doctest vet ./session/tests/once/
doctest test ./session/tests/once/    # expect RED until top-level session package lands
go test ./session/
```

```go
import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"

	"github.com/xhd2015/doctest/session"
)

// Request configures one Once scenario. Leaves set fields via Setup.
type Request struct {
	// SessionID is written into the process env before Once.
	// Empty string means unset DOCTEST_SESSION_ID (missing-session case).
	SessionID string
	// Key is the Once key (may be empty for validation leaves).
	Key string
	// SecondKey, when non-empty, runs a second independent Once (different-keys).
	SecondKey string
	// Mode selects fn behavior: "json-object", "error", "cache-probe".
	Mode string
	// JSONPayload is the object body for success modes (default {"n":1}).
	JSONPayload string
	// CallTwice, when true, invokes Once twice with the same key.
	CallTwice bool
}

// Response captures Once results for Assert.
type Response struct {
	Value       json.RawMessage
	SecondValue json.RawMessage
	Err         error
	SecondErr   error
	FnCalls     int32
	CacheDir    string
	ValueFile   []byte
	ErrorFile   string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	// Control product session id via syscall.Setenv (not t.Setenv/os.Setenv).
	// t.Setenv uses os.LookupEnv which is recorded in go's testlog and would
	// pin the whole workspace package cache to the outer DOCTEST_SESSION_ID.
	// Product Once reads via syscall.Getenv only.
	prev, had := syscall.Getenv(session.DoctestSessionIDEnv)
	if req.SessionID == "" {
		_ = syscall.Unsetenv(session.DoctestSessionIDEnv)
	} else {
		_ = syscall.Setenv(session.DoctestSessionIDEnv, req.SessionID)
	}
	t.Cleanup(func() {
		if had {
			_ = syscall.Setenv(session.DoctestSessionIDEnv, prev)
		} else {
			_ = syscall.Unsetenv(session.DoctestSessionIDEnv)
		}
	})

	resp := &Response{}
	var calls atomic.Int32
	var seenCache string

	fn := func(tb testing.TB, cacheDir string) (json.RawMessage, error) {
		calls.Add(1)
		seenCache = cacheDir
		switch req.Mode {
		case "error":
			return nil, errors.New("boom")
		case "cache-probe":
			marker := filepath.Join(cacheDir, "probe-write")
			if err := os.WriteFile(marker, []byte("ok"), 0o644); err != nil {
				return nil, err
			}
			payload := req.JSONPayload
			if payload == "" {
				payload = `{"path":"` + strings.ReplaceAll(marker, `\`, `\\`) + `"}`
			}
			return json.RawMessage(payload), nil
		default: // "json-object"
			payload := req.JSONPayload
			if payload == "" {
				payload = `{"n":1,"label":"once"}`
			}
			return json.RawMessage(payload), nil
		}
	}

	key := req.Key
	v1, err1 := session.Once(t, key, fn)
	resp.Value = v1
	resp.Err = err1
	resp.CacheDir = seenCache

	if req.CallTwice {
		v2, err2 := session.Once(t, key, fn)
		resp.SecondValue = v2
		resp.SecondErr = err2
	}

	if req.SecondKey != "" {
		vB, errB := session.Once(t, req.SecondKey, func(tb testing.TB, cacheDir string) (json.RawMessage, error) {
			calls.Add(1)
			return json.RawMessage(`{"key":"B"}`), nil
		})
		resp.SecondValue = vB
		resp.SecondErr = errB
	}

	resp.FnCalls = calls.Load()

	// Best-effort disk inspection for layout leaves.
	if seenCache != "" {
		if b, err := os.ReadFile(filepath.Join(seenCache, "value")); err == nil {
			resp.ValueFile = b
		}
		if b, err := os.ReadFile(filepath.Join(seenCache, "error")); err == nil {
			resp.ErrorFile = strings.TrimSpace(string(b))
		}
	}

	// Never fail Run on Once errors — Assert owns expectations.
	return resp, nil
}

// ensureSyscallSees is a harness note: product must use syscall.Getenv.
// Do not read DOCTEST_SESSION_ID via os.Getenv in this tree.
var _ = syscall.Getenv
```
