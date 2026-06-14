## Preconditions
- A project with a doctest root exists but no ASSERT.md leaves.

## Steps
1. Create a project with a doctest root (DOCTEST.md + SETUP.md) but no leaf directories.
2. Run `doctest test <test-dir>`.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	testDir := createTempProject(t, "tests")
	req.Args = append(req.Args, testDir)
	return nil
}
```
