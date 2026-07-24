# Scenario

**Feature**: changing a file under leaf `testdata/` selects that leaf

```
# leaf_a/testdata/input.txt in changed list
FilterByChangedFiles -> [leaf_a]
```

## Steps

1. Create flat two-leaf tree with `leaf_a/testdata/input.txt` present.
2. Set changed path to that testdata file.
3. Run filter policy.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	testdataDir := filepath.Join(fx.TreeDir, "leaf_a", "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testdataDir, "input.txt"), []byte("baseline\n"), 0644); err != nil {
		t.Fatal(err)
	}
	applyPolicyBase(req, fx)
	req.ChangedFiles = []string{treeRel(fx, "leaf_a", "testdata", "input.txt")}
	return nil
}
```
