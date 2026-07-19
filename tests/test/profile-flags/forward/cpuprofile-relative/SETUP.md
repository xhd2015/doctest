# Scenario

**Feature**: relative `-cpuprofile` is abs-resolved against process cwd and appears on the go test line

```
# cwd is WorkDir; relative path becomes absolute on go command line
doctest test -v -cpuprofile profiles/cpu.out <dir>
  -> stderr contains -cpuprofile=<abs-cwd>/profiles/cpu.out
```

## Preconditions
- A dedicated WorkDir is used as the CLI process cwd so abs resolution is deterministic.

## Steps
1. Create a session-scoped work directory.
2. Run `doctest test -v -cpuprofile profiles/cpu.out <fixture>`.

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	wd := filepath.Join(os.TempDir(), "doctest-profile-flags-"+d.DOCTEST_SESSION_ID, "cpu-rel")
	if err := os.MkdirAll(wd, 0o755); err != nil {
		t.Fatalf("mkdir workdir: %v", err)
	}
	req.WorkDir = wd

	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test", "-v",
		"-cpuprofile", "profiles/cpu.out",
		exampleDir,
	}
	return nil
}
```
