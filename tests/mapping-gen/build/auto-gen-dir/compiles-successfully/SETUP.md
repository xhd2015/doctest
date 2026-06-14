## Preconditions
- A project with 1 leaf exists.
- No --gen-dir is specified; build uses a temp directory.

## Steps
1. Create a project with 1 leaf.
2. Run `doctest build <test-dir> -v`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	testDir := createTempProject(t, "tests")
	if err := createDoctestLeaf(filepath.Join(testDir, "simple")); err != nil {
		t.Fatalf("create leaf: %v", err)
	}
	req.Args = append(req.Args, "-v", testDir)
	return nil
}
```
