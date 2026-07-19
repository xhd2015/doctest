# Scenario

**Feature**: absolute `-memprofile` is forwarded as-is on the go test command line

```
doctest test -v -memprofile /abs/session/mem.out <dir>
  -> stderr contains -memprofile=/abs/session/mem.out
```

## Preconditions
- An absolute path under a session temp directory is used for the profile file.

## Steps
1. Build an absolute memprofile path under the session temp dir.
2. Run `doctest test -v -memprofile <abs> <fixture>`.

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(os.TempDir(), "doctest-profile-flags-"+d.DOCTEST_SESSION_ID, "mem-abs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	absMem := filepath.Join(dir, "mem.out")

	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test", "-v",
		"-memprofile", absMem,
		exampleDir,
	}
	return nil
}
```
