# Scenario

**Feature**: `-memprofilerate` is forwarded with the exact value when set, including `0`

```
doctest test -v -memprofilerate 0 -memprofile <abs> <dir>
  -> stderr contains -memprofilerate=0 (or space form) and -memprofile=...
```

## Preconditions
- Rate flags use presence tracking: unset → omit; set → forward exact value including zero.

## Steps
1. Run with `-memprofilerate 0` and an absolute `-memprofile` so the rate is useful to go.

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(os.TempDir(), "doctest-profile-flags-"+d.DOCTEST_SESSION_ID, "memrate")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	absMem := filepath.Join(dir, "mem.out")

	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test", "-v",
		"-memprofilerate", "0",
		"-memprofile", absMem,
		exampleDir,
	}
	return nil
}
```
