# Scenario

**Feature**: `doctest skills update` exit status and stdout

```
optional skill install -> doctest skills update -> {stdout, exit code}
```

## Preconditions

- `req.Bin` and `req.WorkDir` are set by root setup.

## Steps

1. Leaves configure `PreInstalls` and `Args` for the update invocation.

## Context

- Splits on install location (none / project-local / global-only) and update
  scope flags (`--global` or default).
- Registry CLI names match `doctest skill --list` (stable sorted order).

```go
import (
	"strings"
	"testing"
)

func registryCLINames() []string {
	// Must match libdoc/spec registry keys sorted alphabetically (skills update order).
	return []string{
		"analyse-perf",
		"code-spec",
		"design-principle",
		"designer",
		"dev-test",
		"doc-spec",
		"implementer",
		"lint",
		"migrate",
		"output-assert",
		"reproduce",
		"review",
		"review-perf",
		"tdd",
		"tdd-cli-agent",
		"tdd-lite",
	}
}

// skills v0.0.26+ polished batch lines: "<name>  <status>" (+ summary).
// Status tokens may be gray ANSI when color is on; match name + plain status.
func assertNotInstalledLines(t *testing.T, stdout string, names ...string) {
	t.Helper()
	plain := stripANSI(stdout)
	for _, name := range names {
		want := name + "  not installed"
		if !strings.Contains(plain, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
}

func assertUpToDateLine(t *testing.T, stdout, name string) {
	t.Helper()
	want := name + "  up to date"
	if !strings.Contains(stripANSI(stdout), want) {
		t.Fatalf("stdout missing %q:\n%s", want, stdout)
	}
}

func assertNoScopeHint(t *testing.T, stdout string) {
	t.Helper()
	if strings.Contains(stdout, "No installed skills found") {
		t.Fatalf("aggregate scope hint must be removed:\n%s", stdout)
	}
}

// stripANSI removes CSI color sequences so status checks work with color auto/always.
func stripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && (s[j] < 0x40 || s[j] > 0x7e) {
				j++
			}
			if j < len(s) {
				i = j
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"skills", "update"}
	}
	return nil
}
```