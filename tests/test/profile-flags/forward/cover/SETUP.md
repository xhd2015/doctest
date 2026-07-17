# Scenario

**Feature**: `-cover` and relative `-coverprofile` are forwarded; coverprofile path is abs-resolved

```
doctest test -v -cover -coverprofile cov/cover.out <dir>
  -> stderr contains -cover and -coverprofile=<abs>/cov/cover.out
```

## Preconditions
- WorkDir is session-scoped for relative coverprofile resolution.

## Steps
1. Create WorkDir; run with `-cover` and relative `-coverprofile`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	wd := filepath.Join(os.TempDir(), "doctest-profile-flags-"+DOCTEST_SESSION_ID, "cover")
	if err := os.MkdirAll(wd, 0o755); err != nil {
		t.Fatalf("mkdir workdir: %v", err)
	}
	req.WorkDir = wd

	exampleDir := filepath.Join(DOCTEST_ROOT, "testdata", "basic-request-runner")
	req.Args = []string{
		"test", "-v",
		"-cover",
		"-coverprofile", "cov/cover.out",
		exampleDir,
	}
	return nil
}
```
