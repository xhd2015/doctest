package core

import (
	"os"

	"github.com/google/uuid"
)

// DoctestSessionIDEnv is the environment variable set for each doctest test run.
// Generated tests read it via syscall.Getenv into the DOCTEST_SESSION_ID variable.
const DoctestSessionIDEnv = "DOCTEST_SESSION_ID"

// NewDoctestSessionID returns a new session identifier for one doctest test run.
func NewDoctestSessionID() string {
	return uuid.NewString()
}

// DoctestSessionIDForRun returns the session id for the current doctest test
// invocation. It reuses DOCTEST_SESSION_ID from the environment when already set;
// otherwise it generates a new UUID.
func DoctestSessionIDForRun() string {
	if v := os.Getenv(DoctestSessionIDEnv); v != "" {
		return v
	}
	return NewDoctestSessionID()
}