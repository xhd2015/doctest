## Preconditions
- No --gen-dir is specified; tests are generated under the mapping-gen cache root.

## Steps
1. Create a project with 1 leaf.
2. Run `doctest test <test-dir> -v`.

```go
import (
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
