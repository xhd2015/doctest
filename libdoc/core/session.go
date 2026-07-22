package core

import (
	"syscall"

	"github.com/google/uuid"
)

// DoctestSessionIDEnv is the process env key doctest sets on child go test
// processes via cmd.Env (key-replace). Generated test boilerplate reads it via
// syscall.Getenv into the DOCTEST_SESSION_ID variable; harness SETUP/ASSERT code
// must use that variable directly, not os.Getenv.
//
// Product code must not os.Setenv this key — hold it on Options.SessionID for
// the CLI run and pass via ChildEnv / goTestEnv only.
const DoctestSessionIDEnv = "DOCTEST_SESSION_ID"

// NewDoctestSessionID returns a new session identifier for one doctest test run.
func NewDoctestSessionID() string {
	return uuid.NewString()
}

// DoctestSessionIDForRun returns the session id inherited by this process via
// DOCTEST_SESSION_ID (nested go test / suite child), or mints a new UUID when
// unset. Prefer Options.SessionID (set once at CLI Test entry) during a run so
// all go tool children share the same id without process Setenv.
//
// Uses syscall.Getenv so reads are not recorded in the go test cache input log.
func DoctestSessionIDForRun() string {
	v, ok := syscall.Getenv(DoctestSessionIDEnv)
	if ok && v != "" {
		return v
	}
	return NewDoctestSessionID()
}

// SessionIDFromOpts returns opts.SessionID when set; otherwise falls back to
// DoctestSessionIDForRun (inherit or mint). Library callers that do not go
// through runner.Test should still prefer setting opts.SessionID once.
func SessionIDFromOpts(opts Options) string {
	if opts.SessionID != "" {
		return opts.SessionID
	}
	return DoctestSessionIDForRun()
}