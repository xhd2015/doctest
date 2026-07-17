// Package session provides cross-process, session-scoped run-once helpers for
// doctest integration leaves (each leaf runs as its own package/process).
package session

import (
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
var processMemo sync.Map // memoKey -> string

// Once runs fn at most once per (DOCTEST_SESSION_ID, key) on this machine.
//
// cacheDir passed to fn is:
//
//	${UserCacheDir}/doctest/sessions/<session-id>/once-<slug(key)>/
//
// The returned string is written to cacheDir/value and returned to every
// subsequent caller with the same session and key. Use the string as a path,
// URL, or other opaque handle — not as a stand-in for in-memory objects.
//
// Session id is read with syscall.Getenv only (never os.Getenv).
func Once(t testing.TB, key string, fn func(t testing.TB, cacheDir string) (string, error)) (string, error) {
	t.Helper()
	if key == "" {
		return "", errEmptyKey
	}
	sid, ok := syscall.Getenv(DoctestSessionIDEnv)
	if !ok || sid == "" {
		return "", errMissingSession
	}
	slug := slugify(key)
	if slug == "" {
		return "", errEmptySlug
	}

	memoKey := sid + "\x00" + key
	if v, ok := processMemo.Load(memoKey); ok {
		return v.(string), nil
	}

	cacheDir, err := onceDir(sid, slug)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", fmt.Errorf("session.Once: mkdir: %w", err)
	}

	lockPath := filepath.Join(cacheDir, "lock")
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return "", fmt.Errorf("session.Once: open lock: %w", err)
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return "", fmt.Errorf("session.Once: flock: %w", err)
	}
	defer func() { _ = syscall.Flock(int(lock.Fd()), syscall.LOCK_UN) }()

	valuePath := filepath.Join(cacheDir, "value")
	errorPath := filepath.Join(cacheDir, "error")

	if b, err := os.ReadFile(valuePath); err == nil {
		s := strings.TrimSuffix(string(b), "\n")
		if s != "" {
			processMemo.Store(memoKey, s)
			return s, nil
		}
	}
	if b, err := os.ReadFile(errorPath); err == nil && len(b) > 0 {
		return "", errors.New(strings.TrimSpace(string(b)))
	}

	// Starter (or retry after incomplete previous attempt with no value file).
	if fn == nil {
		return "", errors.New("session.Once: fn is nil")
	}
	s, err := fn(t, cacheDir)
	if err != nil {
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		_ = os.Remove(valuePath)
		return "", err
	}
	if s == "" {
		err := errors.New("session.Once: fn returned empty string")
		_ = os.WriteFile(errorPath, []byte(err.Error()+"\n"), 0o644)
		return "", err
	}
	if err := writeFileAtomic(valuePath, []byte(s+"\n")); err != nil {
		return "", fmt.Errorf("session.Once: write value: %w", err)
	}
	_ = os.Remove(errorPath)
	processMemo.Store(memoKey, s)
	return s, nil
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
