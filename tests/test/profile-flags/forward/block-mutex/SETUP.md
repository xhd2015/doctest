# Scenario

**Feature**: block/mutex profile paths and rates are forwarded to go test

```
doctest test -v \
  -blockprofile <abs>/block.out -blockprofilerate 1 \
  -mutexprofile <abs>/mutex.out -mutexprofilefraction 1 \
  <dir>
  -> stderr contains all four flags with exact values
```

## Preconditions
- Absolute profile paths under a session temp directory.
- Rates set to non-default positive integers so presence is unambiguous.

## Steps
1. Create session profile dir; run with block and mutex profile + rate flags.

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(os.TempDir(), "doctest-profile-flags-"+d.DOCTEST_SESSION_ID, "block-mutex")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	block := filepath.Join(dir, "block.out")
	mutex := filepath.Join(dir, "mutex.out")

	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test", "-v",
		"-blockprofile", block,
		"-blockprofilerate", "1",
		"-mutexprofile", mutex,
		"-mutexprofilefraction", "1",
		exampleDir,
	}
	return nil
}
```
