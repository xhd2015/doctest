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
	return []string{
		"code-spec",
		"designer",
		"doc-spec",
		"implementer",
		"output-assert",
		"reproduce",
		"review",
		"tdd",
		"tdd-lite",
	}
}

func assertNotInstalledLines(t *testing.T, stdout string, names ...string) {
	t.Helper()
	for _, name := range names {
		want := "skill not installed: " + name
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
}

func assertNoScopeHint(t *testing.T, stdout string) {
	t.Helper()
	if strings.Contains(stdout, "No installed skills found") {
		t.Fatalf("aggregate scope hint must be removed:\n%s", stdout)
	}
}

func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 {
		req.Args = []string{"skills", "update"}
	}
	return nil
}
```