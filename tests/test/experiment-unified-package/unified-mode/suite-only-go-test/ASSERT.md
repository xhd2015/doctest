---
label: heavy
---

## Expected

- Suite run succeeds.
- A displayed `go test` line is present.
- Package path args list has **exactly one** entry.
- That entry contains `suite` (e.g. `./suite` or `./tree/suite`).
- Package args do **not** include both leaf packages `./a` and `./b` (or `./…/a` and `./…/b` as separate multi-package run).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("unified-mode RunTest failed: %s\nstdout:\n%s\nstderr:\n%s",
			resp.RunErr, resp.Stdout, resp.Stderr)
	}
	if resp.GoTestDisplayLine == "" {
		t.Fatalf("no go test display line found in stdout/stderr:\nstdout:\n%s\nstderr:\n%s",
			resp.Stdout, resp.Stderr)
	}
	pkgs := resp.GoTestPackageArgs
	if len(pkgs) != 1 {
		t.Fatalf("unified mode wants exactly 1 go test package (suite), got %d args=%v line=%q",
			len(pkgs), pkgs, resp.GoTestDisplayLine)
	}
	if !strings.Contains(pkgs[0], "suite") {
		t.Fatalf("single package arg must refer to suite, got %q line=%q", pkgs[0], resp.GoTestDisplayLine)
	}
	// Hard fail if multi-leaf package invocation shape leaked through.
	joined := strings.Join(pkgs, " ")
	if (strings.Contains(joined, "/a") || strings.HasSuffix(joined, "a") || strings.Contains(joined, "./a")) &&
		(strings.Contains(joined, "/b") || strings.HasSuffix(joined, "b") || strings.Contains(joined, "./b")) {
		t.Fatalf("go test still multi-package leaf shape: %v line=%q", pkgs, resp.GoTestDisplayLine)
	}
}
```
