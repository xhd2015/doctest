---
explanation: subset does not heal out-of-scope missing; tree-b run does
---

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	requireExit0(t, "full", resp.FullExitCode, resp.FullStderr)
	requireExit0(t, "subset tree-a", resp.SubsetExitCode, resp.SubsetStderr)
	requireExit0(t, "tree-b heal", resp.ThirdExitCode, resp.ThirdStderr)

	if !manifestHasExactRel(resp.ManifestAfterFull, "tree-b/b1/leaf.go") {
		t.Fatalf("full manifest missing tree-b/b1/leaf.go:\n%s", resp.ManifestAfterFull)
	}
	// Subset must not restore out-of-scope file.
	if resp.DeletedGenExistsAfter2 {
		t.Fatalf("subset healed out-of-scope file (must stay missing): %s", resp.DeletedGenAbs)
	}
	// Ledger still tracks sibling.
	if !manifestHasExactRel(resp.ManifestAfterSubset, "tree-b/b1/leaf.go") {
		t.Fatalf("subset must not drop sibling ledger entry:\n%s", resp.ManifestAfterSubset)
	}
	for _, line := range strings.Split(resp.SubsetStderr, "\n") {
		if strings.Contains(line, "tree-b") && strings.Contains(line, "leaf.go") && strings.Contains(line, "# new") {
			t.Fatalf("subset result mentions tree-b leaf as new: %s", line)
		}
	}

	// run3 heals.
	if !resp.DeletedGenExistsAfter3 {
		t.Fatalf("tree-b run should recreate %s", resp.DeletedGenAbs)
	}
	if !genPlanHasTag(resp.ThirdStderr, "leaf.go", "new") {
		if parseNewFromGenPlan(resp.ThirdStderr) < 1 {
			t.Fatalf("tree-b heal expected # new; stderr:\n%s", resp.ThirdStderr)
		}
	}
}
```
