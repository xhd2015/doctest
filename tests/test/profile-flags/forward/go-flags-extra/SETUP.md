# Scenario

**Feature**: extra go test allowlist flags are forwarded on the verbose go command line

```
doctest test -v \
  -covermode atomic -coverpkg example.com/mod/... \
  -short -failfast -parallel 2 -shuffle on \
  -tags integration -gcflags all=-N -ldflags -X=main.v=1 \
  -race \
  <dir>
  -> stderr go line contains each flag
```

## Preconditions
- Single-package fixture so multi-package coverprofile rules do not apply.
- Nested go test; use e2e when full-integration.

## Steps
1. Run nested doctest test with the allowlist flags against basic-request-runner.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	exampleDir := filepath.Join(d.DOCTEST_ROOT, "testdata", "basic-request-runner")
	if _, err := os.Stat(exampleDir); err != nil {
		t.Fatalf("fixture %s: %v", exampleDir, err)
	}
	req.Args = []string{
		"test", "-v",
		"-covermode", "atomic",
		"-coverpkg", "example.com/mod/...",
		"-short",
		"-failfast",
		"-parallel", "2",
		"-shuffle", "on",
		"-tags", "integration",
		"-gcflags", "all=-N",
		"-ldflags", "-s",
		"-race",
		exampleDir,
	}
	return nil
}
```
