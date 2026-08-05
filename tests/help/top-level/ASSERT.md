## Expected
- The command succeeds.
- stdout contains top-level usage.
- stdout lists the major subcommands.

## Exit Code
- Exit code 0.

```go
import (
    "regexp"
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
    }
	for _, want := range []string{"Usage: doctest", "agent", "vet", "design", "build", "test", "skill"} {
        if !strings.Contains(resp.Stdout, want) {
            t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
        }
    }
	// list must appear as its own command entry (not skill --list / --list-sessions).
	// This is the RED→GREEN signal for adding doctest list to top-level usage.
	if !regexp.MustCompile(`(?m)^[ \t]+list\b`).MatchString(resp.Stdout) {
		t.Fatalf("stdout missing list command entry:\n%s", resp.Stdout)
	}
}

```

