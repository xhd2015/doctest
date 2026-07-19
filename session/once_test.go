package session

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
)

func setSession(t *testing.T, id string) {
	t.Helper()
	// Isolate Once disk cache so repeated local `go test` runs with fixed
	// session ids do not skip fn (CI is always clean; local often is not).
	t.Setenv(DoctestCacheHomeEnv, t.TempDir())
	t.Setenv(DoctestSessionIDEnv, id)
	if v, ok := syscall.Getenv(DoctestSessionIDEnv); !ok || v != id {
		t.Fatalf("syscall.Getenv(%s)=%q,%v want %q", DoctestSessionIDEnv, v, ok, id)
	}
}

func TestOnceMissingSession(t *testing.T) {
	t.Setenv(DoctestSessionIDEnv, "")
	processMemo = sync.Map{}
	_, err := Once(t, "k", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		return json.RawMessage(`{"ok":true}`), nil
	})
	if err == nil || !errors.Is(err, errMissingSession) && !stringsContains(err.Error(), "DOCTEST_SESSION_ID") {
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, errMissingSession) {
			// still ok if wrapped
			if err.Error() != errMissingSession.Error() {
				t.Fatalf("err=%v", err)
			}
		}
	}
}

func stringsContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && (func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})()))
}

func TestOnceEmptyKey(t *testing.T) {
	setSession(t, "sess-empty-key")
	processMemo = sync.Map{}
	_, err := Once(t, "", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		return json.RawMessage(`{}`), nil
	})
	if !errors.Is(err, errEmptyKey) {
		t.Fatalf("err=%v", err)
	}
}

func TestOnceRunsFnOncePerSessionKey(t *testing.T) {
	setSession(t, "sess-once-json-1")
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

	v1, err := Once(t, "cli", fn)
	if err != nil {
		t.Fatal(err)
	}
	v2, err := Once(t, "cli", fn)
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
	setSession(t, "sess-keys-json")
	processMemo = sync.Map{}
	var a, b atomic.Int32
	va, err := Once(t, "a", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		a.Add(1)
		return json.RawMessage(`{"k":"A"}`), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	vb, err := Once(t, "b", func(t testing.TB, cacheDir string) (json.RawMessage, error) {
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
	setSession(t, "sess-err-json")
	processMemo = sync.Map{}
	var calls atomic.Int32
	fn := func(t testing.TB, cacheDir string) (json.RawMessage, error) {
		calls.Add(1)
		return nil, errors.New("boom")
	}
	_, err := Once(t, "fail", fn)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("err=%v", err)
	}
	_, err = Once(t, "fail", fn)
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
