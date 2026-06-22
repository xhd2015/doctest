package core

import (
	"syscall"
	"testing"
)

func TestDoctestSessionIDForRunReusesEnv(t *testing.T) {
	t.Setenv(DoctestSessionIDEnv, "preset-session")
	if got := DoctestSessionIDForRun(); got != "preset-session" {
		t.Fatalf("DoctestSessionIDForRun() = %q, want preset-session", got)
	}
}

func TestDoctestSessionIDForRunGeneratesWhenUnset(t *testing.T) {
	t.Setenv(DoctestSessionIDEnv, "")
	if id := DoctestSessionIDForRun(); id == "" {
		t.Fatal("expected non-empty generated session id")
	}
}

func TestDoctestSessionIDForRunUsesSyscallGetenv(t *testing.T) {
	t.Setenv(DoctestSessionIDEnv, "probe-session")
	v, ok := syscall.Getenv(DoctestSessionIDEnv)
	if !ok || v != "probe-session" {
		t.Fatalf("syscall.Getenv(%q) = %q, %v", DoctestSessionIDEnv, v, ok)
	}
	if got := DoctestSessionIDForRun(); got != "probe-session" {
		t.Fatalf("DoctestSessionIDForRun() = %q, want probe-session", got)
	}
}