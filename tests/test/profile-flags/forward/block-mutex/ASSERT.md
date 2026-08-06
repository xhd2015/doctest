---
explanation: nested go test with block/mutex profiles
---

## Expected
- Exit code 0.
- stderr contains `-blockprofile=…`, `-blockprofilerate=1`, `-mutexprofile=…`, `-mutexprofilefraction=1`.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func flagOnLine(stderr, name, value string) bool {
	eq := name + "=" + value
	sp := name + " " + value
	return strings.Contains(stderr, eq) || strings.Contains(stderr, sp)
}

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	// Recover absolute paths from Args.
	var block, mutex string
	for i, a := range req.Args {
		if a == "-blockprofile" && i+1 < len(req.Args) {
			block = req.Args[i+1]
		}
		if a == "-mutexprofile" && i+1 < len(req.Args) {
			mutex = req.Args[i+1]
		}
	}
	if block == "" || mutex == "" {
		t.Fatalf("missing profile paths in req.Args: %v", req.Args)
	}
	// Normalize in case of relative (should already be abs).
	if !filepath.IsAbs(block) {
		t.Fatalf("block path should be absolute in this leaf: %s", block)
	}

	checks := []struct {
		name, value string
	}{
		{"-blockprofile", block},
		{"-blockprofilerate", "1"},
		{"-mutexprofile", mutex},
		{"-mutexprofilefraction", "1"},
	}
	for _, c := range checks {
		if !flagOnLine(resp.Stderr, c.name, c.value) {
			t.Fatalf("expected %s=%s on go command line, stderr:\n%s", c.name, c.value, resp.Stderr)
		}
	}
}
```
