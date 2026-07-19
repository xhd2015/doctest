# Scenario

**Feature**: relative `-outputdir` is abs-resolved and forwarded to go test

```
doctest test -v -outputdir outprof <dir>
  -> stderr contains -outputdir=<abs-cwd>/outprof
```

## Preconditions
- WorkDir is session-scoped.

## Steps
1. Create WorkDir; run with `-outputdir outprof`.

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	wd := filepath.Join(os.TempDir(), "doctest-profile-flags-"+d.DOCTEST_SESSION_ID, "outputdir")
	if err := os.MkdirAll(wd, 0o755); err != nil {
		t.Fatalf("mkdir workdir: %v", err)
	}
	req.WorkDir = wd

	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test", "-v",
		"-outputdir", "outprof",
		// Pair with a profile so outputdir is meaningful to go; still assert command line only.
		"-cpuprofile", "cpu.out",
		exampleDir,
	}
	return nil
}
```
