# Scenario

**Feature**: `-coverprofile` writes a coverage profile file after a successful run

```
doctest test -cover -coverprofile <session>/cover.out <fixture>
  -> exit 0 -> cover.out exists
```

## Preconditions
- Absolute coverprofile path under session temp; single-package fixture.

## Steps
1. Create session path; run with `-cover` and `-coverprofile`.

```go
import (
"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(os.TempDir(), "doctest-profile-flags-"+d.DOCTEST_SESSION_ID, "side-cover")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	covPath := filepath.Join(dir, "cover.out")
	_ = os.Remove(covPath)

	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test",
		"-cover",
		"-coverprofile", covPath,
		exampleDir,
	}
	return nil
}
```
