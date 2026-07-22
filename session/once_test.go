package session

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
)

func TestOnceMissingSession(t *testing.T) {
	// Production Once reads process env; with empty env it fails without Setenv.
	// Isolate from suite/parent env if present.
	processMemo = sync.Map{}
	// Clear via OnceSession empty is separate; test env path:
	// If parent already set DOCTEST_SESSION_ID, Once would succeed — use OnceSession
	// for empty-sid and keep Once env test only when we can unset without t.Setenv.
	// Prefer OnceSession for Parallel-safe missing-sid check:
	_, err := OnceSession(t, "", "k", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		return json.RawMessage(`{"ok":true}`), nil
	})
	if err == nil {
		t.Fatal("expected error for empty session id")
	}
	if !errors.Is(err, errEmptySession) && err.Error() != errEmptySession.Error() {
		if !stringsContains(err.Error(), "session") {
			t.Fatalf("err=%v", err)
		}
	}
}

func TestOnceEnvMissingSession(t *testing.T) {
	// Cover production Once(getenv) path: only when env is actually empty.
	// Skip if outer process already injected a session id (e.g. under doctest test).
	processMemo = sync.Map{}
	if v, ok := lookupSessionEnv(); ok && v != "" {
		t.Skip("DOCTEST_SESSION_ID already set in process; skip getenv-missing unit")
	}
	_, err := Once(t, "k", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		return json.RawMessage(`{"ok":true}`), nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, errMissingSession) && err.Error() != errMissingSession.Error() {
		if !stringsContains(err.Error(), "DOCTEST_SESSION_ID") {
			t.Fatalf("err=%v", err)
		}
	}
}

func lookupSessionEnv() (string, bool) {
	// Use syscall.Getenv via Once's contract — import syscall in test.
	return lookupSessionEnvImpl()
}

func stringsContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}

func TestOnceEmptyKey(t *testing.T) {
	processMemo = sync.Map{}
	_, err := OnceSession(t, "sess-empty-key", "", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		return json.RawMessage(`{}`), nil
	})
	if !errors.Is(err, errEmptyKey) {
		t.Fatalf("err=%v", err)
	}
}

func TestOnceRunsFnOncePerSessionKey(t *testing.T) {
	processMemo = sync.Map{}
	var calls atomic.Int32
	fn := func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		calls.Add(1)
		if cacheDir == "" {
			t.Fatal("empty cacheDir")
		}
		p := filepath.Join(cacheDir, "marker")
		if err := os.WriteFile(p, []byte("1"), 0o644); err != nil {
			return nil, err
		}
		return json.Marshal(map[string]string{"path": p})
	}

	v1, err := OnceSession(t, "sess-once-json-1", "cli", fn)
	if err != nil {
		t.Fatal(err)
	}
	v2, err := OnceSession(t, "sess-once-json-1", "cli", fn)
	if err != nil {
		t.Fatal(err)
	}
	if string(v1) != string(v2) {
		t.Fatalf("values differ: %s vs %s", v1, v2)
	}
	if calls.Load() != 1 {
		t.Fatalf("fn calls=%d want 1", calls.Load())
	}
	var h struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(v1, &h); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(h.Path); err != nil {
		t.Fatalf("marker: %v", err)
	}
}

func TestOnceDifferentKeysIndependent(t *testing.T) {
	processMemo = sync.Map{}
	var a, b atomic.Int32
	va, err := OnceSession(t, "sess-keys-json", "a", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		a.Add(1)
		return json.RawMessage(`{"k":"A"}`), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	vb, err := OnceSession(t, "sess-keys-json", "b", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		b.Add(1)
		return json.RawMessage(`{"k":"B"}`), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(va) != `{"k":"A"}` || string(vb) != `{"k":"B"}` {
		t.Fatalf("got %s %s", va, vb)
	}
	if a.Load() != 1 || b.Load() != 1 {
		t.Fatalf("calls a=%d b=%d", a.Load(), b.Load())
	}
}

func TestOncePropagatesError(t *testing.T) {
	processMemo = sync.Map{}
	var calls atomic.Int32
	fn := func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		calls.Add(1)
		return nil, errors.New("boom")
	}
	_, err := OnceSession(t, "sess-err-json", "fail", fn)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("err=%v", err)
	}
	_, err = OnceSession(t, "sess-err-json", "fail", fn)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("second err=%v", err)
	}
	if calls.Load() != 1 {
		t.Fatalf("fn calls=%d want 1", calls.Load())
	}
}

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"cli":       "cli",
		"go-binary": "go-binary",
		"a/b":       "a-b",
		"foo:bar":   "foo-bar",
	}
	for in, want := range cases {
		if got := slugify(in); got != want {
			t.Errorf("slugify(%q)=%q want %q", in, got, want)
		}
	}
}
