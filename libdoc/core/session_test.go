package core

import (
	"os"
	"testing"
)

func TestDoctestSessionIDForRunReusesEnv(t *testing.T) {
	t.Setenv(DoctestSessionIDEnv, "preset-session")
	if got := DoctestSessionIDForRun(); got != "preset-session" {
		t.Fatalf("DoctestSessionIDForRun() = %q, want preset-session", got)
	}
}

func TestDoctestSessionIDForRunGeneratesWhenUnset(t *testing.T) {
	_ = os.Unsetenv(DoctestSessionIDEnv)
	if id := DoctestSessionIDForRun(); id == "" {
		t.Fatal("expected non-empty generated session id")
	}
}