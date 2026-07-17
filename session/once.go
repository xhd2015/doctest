// Package session provides cross-process, session-scoped run-once helpers for
// doctest integration leaves (each leaf runs as its own package/process).
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
	errEmptyKey       = errors.New("session.Once: key is empty")
	errEmptySlug      = errors.New("session.Once: key slugifies to empty string")
)

// processMemo avoids re-reading disk within one process after a successful Once.
var processMemo sync.Map // memoKey -> json.RawMessage

// Once runs fn at most once per (DOCTEST_SESSION_ID, key) on this machine.
//
// cacheDir passed to fn is:
//
//	${UserCacheDir}/doctest/sessions/<session-id>/once-<slug(key)>/
//
// The returned json.RawMessage is written to cacheDir/value as raw JSON bytes
// and returned to every subsequent caller with the same session and key.
//
// Session id is read with syscall.Getenv only (never os.Getenv).
func Once(t testing.TB, key string, fn func(t testing.TB, cacheDir string) (json.RawMessage, error)) (json.RawMessage, error) {
	t.Helper()
	if key == "" {
		return nil, errEmptyKey
	}
	sid, ok := syscall.Getenv(DoctestSessionIDEnv)
	if !ok || sid == "" {
		return nil, errMissingSession
	}
	slug := slugify(key)
	if slug == "" {
		return nil, errEmptySlug
	}

	memoKey := sid + "\x00" + key
	if v, ok := processMemo.Load(memoKey); ok {
		return append(json.RawMessage(nil), v.(json.RawMessage)...), nil
	}

	cacheDir, err := onceDir(sid, slug)
	if err != nil {
		return nil, err
	}
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

	valuePath := filepath.Join(cacheDir, "value")
	errorPath := filepath.Join(cacheDir, "error")

	if b, err := os.ReadFile(valuePath); err == nil && len(b) > 0 && json.Valid(b) {
		raw := json.RawMessage(append([]byte(nil), b...))
		processMemo.Store(memoKey, raw)
		return append(json.RawMessage(nil), raw...), nil
	}
	if b, err := os.ReadFile(errorPath); err == nil && len(b) > 0 {
		return nil, errors.New(strings.TrimSpace(string(b)))
	}

	// Starter (or retry after incomplete previous attempt with no value file).
	if fn == nil {
		return nil, errors.New("session.Once: fn is nil")
	}
	raw, err := fn(t, cacheDir)
	if err != nil {
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		_ = os.Remove(valuePath)
		return nil, err
	}
	if len(raw) == 0 {
		err := errors.New("session.Once: fn returned empty json.RawMessage")
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		return nil, err
	}
	if !json.Valid(raw) {
		err := errors.New("session.Once: fn returned invalid JSON")
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		return nil, err
	}
	// Persist raw JSON bytes exactly (no trailing newline injection).
	stored := append([]byte(nil), raw...)
	if err := writeFileAtomic(valuePath, stored); err != nil {
		return nil, fmt.Errorf("session.Once: write value: %w", err)
	}
	_ = os.Remove(errorPath)
	out := json.RawMessage(stored)
	processMemo.Store(memoKey, out)
	return append(json.RawMessage(nil), out...), nil
}

func onceDir(sessionID, slug string) (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}
	// Session id is UUID-like; still scrub path separators.
	sid := slugify(sessionID)
	if sid == "" {
		return "", errors.New("session.Once: session id slugifies to empty")
	}
	return filepath.Join(base, "doctest", "sessions", sid, "once-"+slug), nil
}

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
