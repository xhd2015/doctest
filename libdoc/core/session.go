package core

import (
	"syscall"

	"github.com/google/uuid"
)

// DoctestSessionIDEnv is the environment variable set for each doctest test run.
// DoctestSessionIDEnv is the process env key doctest test sets for child go test
// processes. Generated test boilerplate reads it via syscall.Getenv into the
// DOCTEST_SESSION_ID variable; harness SETUP/ASSERT code must use that variable
// directly, not os.Getenv.
const DoctestSessionIDEnv = "DOCTEST_SESSION_ID"

// NewDoctestSessionID returns a new session identifier for one doctest test run.
func NewDoctestSessionID() string {
	return uuid.NewString()
}

// DoctestSessionIDForRun returns the session id for the current doctest test
// invocation. It reuses DOCTEST_SESSION_ID from the environment when already set;
// otherwise it generates a new UUID. Uses syscall.Getenv so reads are not
// recorded in the go test cache input log.
func DoctestSessionIDForRun() string {
	v, ok := syscall.Getenv(DoctestSessionIDEnv)
	if ok && v != "" {
		return v
	}
	return NewDoctestSessionID()
}