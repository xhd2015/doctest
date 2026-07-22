package core

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

// sessionHelperEnv is set on re-exec children that exercise DoctestSessionIDForRun
// without t.Setenv (process env writes are forbidden in unit tests).
const sessionHelperEnv = "DOCTEST_CORE_SESSION_HELPER"

func TestNewDoctestSessionIDNonEmpty(t *testing.T) {
	id := NewDoctestSessionID()
	if id == "" {
		t.Fatal("expected non-empty session id")
	}
	if id2 := NewDoctestSessionID(); id2 == "" || id2 == id {
		// UUID collision is astronomically rare; empty is the real failure mode.
		if id2 == "" {
			t.Fatal("expected non-empty second session id")
		}
	}
}

func TestSessionIDFromOptsPrefersOpts(t *testing.T) {
	got := SessionIDFromOpts(Options{SessionID: "opts-sid"})
	if got != "opts-sid" {
		t.Fatalf("SessionIDFromOpts = %q, want opts-sid", got)
	}
}

func TestDoctestSessionIDForRunReusesEnv(t *testing.T) {
	if os.Getenv(sessionHelperEnv) == "reuse" {
		if got := DoctestSessionIDForRun(); got != "preset-session" {
			t.Fatalf("DoctestSessionIDForRun() = %q, want preset-session", got)
		}
		v, ok := syscall.Getenv(DoctestSessionIDEnv)
		if !ok || v != "preset-session" {
			t.Fatalf("syscall.Getenv(%q) = %q, %v", DoctestSessionIDEnv, v, ok)
		}
		return
	}
	runSessionHelper(t, "reuse", []string{DoctestSessionIDEnv + "=preset-session"})
}

func TestDoctestSessionIDForRunGeneratesWhenUnset(t *testing.T) {
	if os.Getenv(sessionHelperEnv) == "generate" {
		if id := DoctestSessionIDForRun(); id == "" {
			t.Fatal("expected non-empty generated session id")
		}
		return
	}
	// Child env strips DOCTEST_SESSION_ID so ForRun must mint.
	runSessionHelper(t, "generate", nil)
}

func runSessionHelper(t *testing.T, mode string, extra []string) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$", "-test.v=false")
	cmd.Env = childEnvWithoutSession(extra...)
	cmd.Env = append(cmd.Env, sessionHelperEnv+"="+mode)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("helper %s: %v\n%s", mode, err, out)
	}
}

// childEnvWithoutSession copies process env, strips DOCTEST_SESSION_ID, then
// applies extra KEY=val entries (key-replace).
func childEnvWithoutSession(extra ...string) []string {
	base := make([]string, 0, len(os.Environ()))
	for _, e := range os.Environ() {
		k, _, _ := strings.Cut(e, "=")
		if k == DoctestSessionIDEnv || k == sessionHelperEnv {
			continue
		}
		base = append(base, e)
	}
	return ChildEnv(base, extra...)
}
