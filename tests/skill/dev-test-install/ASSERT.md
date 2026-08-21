## Expected

- Installation succeeds.
- The installed `SKILL.md` exactly matches the embedded registry content,
  including `metadata.version`.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/spec"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr: %s", resp.ExitCode, resp.Stderr)
	}
	installed, readErr := os.ReadFile(filepath.Join(req.InstallDir, "SKILL.md"))
	if readErr != nil {
		t.Fatalf("read installed SKILL.md: %v", readErr)
	}
	want, contentErr := spec.Content("dev-test")
	if contentErr != nil {
		t.Fatalf("registry content: %v", contentErr)
	}
	if string(installed) != want {
		t.Fatalf("installed content differs from registry")
	}
	if !strings.Contains(string(installed), "metadata:\n  version: \"0.1.0\"") {
		t.Fatalf("installed SKILL.md missing version metadata:\n%s", installed)
	}
}
```
