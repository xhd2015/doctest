# Scenario

**Feature**: vet selection returns only the changed leaf SETUP (not sibling)

```
# leaf_b SETUP anti-pattern content + listed as changed
ChangedDoctestMarkdownFiles -> [leaf_b/SETUP.md] only
```

## Steps

1. Create two-leaf tree; write anti-pattern content into `leaf_b/SETUP.md`.
2. Set `ChangedFiles` to that SETUP only.
3. Assert markdown path list excludes `leaf_a`.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	fx := createVetFlatTwoLeafTree(t)
	setupPath := filepath.Join(fx.TreeDir, "leaf_b", "SETUP.md")
	if err := os.WriteFile(setupPath, readAntiPatternSetup(t, d), 0644); err != nil {
		t.Fatal(err)
	}
	applyPolicyBase(req, fx)
	req.Policy = PolicyVetMD
	req.ChangedFiles = []string{treeRel(fx, "leaf_b", "SETUP.md")}
	return nil
}
```
