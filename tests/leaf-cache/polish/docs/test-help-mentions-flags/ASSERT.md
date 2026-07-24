## Expected

- Exit 0.
- Help text (stdout or stderr) contains `-a` and wipe/hard-force wording.
- Help must not mention removed `--no-leaf-cache`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("help exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	out := resp.Stdout + "\n" + resp.Stderr
	if !strings.Contains(out, "-a") {
		t.Fatalf("help missing -a:\n%s", out)
	}
	if strings.Contains(out, "--no-leaf-cache") {
		t.Fatalf("help must not document removed --no-leaf-cache:\n%s", out)
	}
	if !strings.Contains(out, "wipe") && !strings.Contains(out, "Hard force") {
		t.Fatalf("help -a should describe hard force / wipe:\n%s", out)
	}
}
```
