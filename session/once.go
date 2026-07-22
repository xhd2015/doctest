// Package session provides session-scoped run-once helpers for doctest leaves.
//
// Disk scratch for Once lives under t.TempDir() so file opens are not Go
// testcache inputs (unlike UserCacheDir/doctest/sessions/<sid>/…, which bust
// package cache when DOCTEST_SESSION_ID changes each CLI run).
// Cross-leaf sharing within one go test process is via processMemo; durable
// cross-process artifacts (e.g. testbin binary path) should live outside Once's
// cacheDir (see libdoc/testbin).
package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"unicode"
)

// DoctestSessionIDEnv is the process env key doctest test sets for child go test
// processes. Read only via syscall.Getenv so Go's test result cache is not
// keyed on this value (os.Getenv is recorded in the testlog).
const DoctestSessionIDEnv = "DOCTEST_SESSION_ID"

var (
	errMissingSession = errors.New("session.Once: DOCTEST_SESSION_ID is not set (syscall.Getenv)")
	errEmptySession   = errors.New("session.OnceSession: session id is empty")
	errEmptyKey       = errors.New("session.Once: key is empty")
	errEmptySlug      = errors.New("session.Once: key slugifies to empty string")
)

// OnceFn is the callback passed to Once / OnceSession.
type OnceFn func(t testing.TB, cacheDir string) (json.RawMessage, error)

type onceMemo struct {
	raw json.RawMessage
	err string // non-empty => remembered failure
}

// processMemo avoids re-running fn within one process after a successful or
// failed Once for the same (session id, key).
var processMemo sync.Map // memoKey -> onceMemo

// Once runs fn at most once per (DOCTEST_SESSION_ID, key) within this process
// (processMemo). Session id is read with syscall.Getenv only (never os.Getenv).
//
// Production path: the suite go test process has DOCTEST_SESSION_ID set once via
// child cmd.Env. Do not mutate process env from concurrent leaves — use
// OnceSession with an explicit id in Parallel-safe tests.
//
// cacheDir passed to fn is:
//
//	t.TempDir()/session-once/<slug(key)>/
func Once(t testing.TB, key string, fn OnceFn) (json.RawMessage, error) {
	t.Helper()
	sid, ok := syscall.Getenv(DoctestSessionIDEnv)
	if !ok || sid == "" {
		return nil, errMissingSession
	}
	return onceImpl(t, sid, key, fn)
}

// OnceSession is Parallel-safe: session id is explicit and process env is not read
// or written. Use this from concurrent leaves / unit tests that need an isolated sid.
// Empty sessionID returns errEmptySession without running fn.
func OnceSession(t testing.TB, sessionID, key string, fn OnceFn) (json.RawMessage, error) {
	t.Helper()
	if sessionID == "" {
		return nil, errEmptySession
	}
	return onceImpl(t, sessionID, key, fn)
}

func onceImpl(t testing.TB, sid, key string, fn OnceFn) (json.RawMessage, error) {
	t.Helper()
	if key == "" {
		return nil, errEmptyKey
	}
	slug := slugify(key)
	if slug == "" {
		return nil, errEmptySlug
	}

	memoKey := sid + "\x00" + key
	if v, ok := processMemo.Load(memoKey); ok {
		return memoResult(v.(onceMemo))
	}

	// TempDir is excluded from go testcache inputs; do not use UserCacheDir
	// sessions/<sid>/ (that path changes every CLI session and busts cache).
	cacheDir := filepath.Join(t.TempDir(), "session-once", slug)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("session.Once: mkdir: %w", err)
	}

	lockPath := filepath.Join(cacheDir, "lock")
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("session.Once: open lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return nil, fmt.Errorf("session.Once: flock: %w", err)
	}
	defer func() { _ = syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) }()

	// Another caller may have filled processMemo while we waited on the lock.
	if v, ok := processMemo.Load(memoKey); ok {
		return memoResult(v.(onceMemo))
	}

	valuePath := filepath.Join(cacheDir, "value")
	errorPath := filepath.Join(cacheDir, "error")

	if b, err := os.ReadFile(valuePath); err == nil && len(b) > 0 && json.Valid(b) {
		raw := json.RawMessage(append([]byte(nil), b...))
		processMemo.Store(memoKey, onceMemo{raw: raw})
		return append(json.RawMessage(nil), raw...), nil
	}
	if b, err := os.ReadFile(errorPath); err == nil && len(b) > 0 {
		msg := strings.TrimSpace(string(b))
		processMemo.Store(memoKey, onceMemo{err: msg})
		return nil, errors.New(msg)
	}

	if fn == nil {
		return nil, errors.New("session.Once: fn is nil")
	}
	raw, err := fn(t, cacheDir)
	if err != nil {
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		_ = os.Remove(valuePath)
		processMemo.Store(memoKey, onceMemo{err: err.Error()})
		return nil, err
	}
	if len(raw) == 0 {
		err := errors.New("session.Once: fn returned empty json.RawMessage")
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		processMemo.Store(memoKey, onceMemo{err: err.Error()})
		return nil, err
	}
	if !json.Valid(raw) {
		err := errors.New("session.Once: fn returned invalid JSON")
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		processMemo.Store(memoKey, onceMemo{err: err.Error()})
		return nil, err
	}
	stored := append([]byte(nil), raw...)
	if err := writeFileAtomic(valuePath, stored); err != nil {
		return nil, fmt.Errorf("session.Once: write value: %w", err)
	}
	_ = os.Remove(errorPath)
	out := json.RawMessage(stored)
	processMemo.Store(memoKey, onceMemo{raw: out})
	return append(json.RawMessage(nil), out...), nil
}

func memoResult(m onceMemo) (json.RawMessage, error) {
	if m.err != "" {
		return nil, errors.New(m.err)
	}
	return append(json.RawMessage(nil), m.raw...), nil
}

// DoctestCacheHomeEnv overrides durable cache roots when set.
// Once no longer stores under this path (uses t.TempDir); kept for API
// compatibility with callers that still reference the name.
const DoctestCacheHomeEnv = "DOCTEST_CACHE_HOME"

// slugify keeps [a-zA-Z0-9._-], maps other runes to '-', collapses repeats.
func slugify(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevDash := false
	for _, r := range s {
		ok := unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '-'
		if ok {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash && b.Len() > 0 {
			b.WriteByte('-')
			prevDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	return out
}

func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".value-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
