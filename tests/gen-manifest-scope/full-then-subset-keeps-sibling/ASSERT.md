---
label: heavy
explanation: multi-tree then subset; sibling stays in doctest.gen-manifest
---

## Expected

- Full run creates `doctest.gen-manifest` listing both `tree-a/` and `tree-b/` paths.
- Subset run (`tree-a` only) **does not remove** `tree-b/` entries from the
  manifest (ledger is not rewritten to the subset desired set).
- Gen dir still contains a `tree-b/` package directory after the subset run.
- Both CLI invocations complete prepare successfully (exit 0 with bypass-go-test,
  or soft non-zero only if prepare fails — require exit 0).

## Errors

- Fail if full manifest lacks either tree.
- Fail if subset shrinks away all `tree-b/` ledger lines.
- Fail if sibling gen dir disappeared (incorrect full-gen prune on subset).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run error: %v\nfull stderr:\n%s\nsubset stderr:\n%s", err, resp.FullStderr, resp.SubsetStderr)
	}
	if resp.FullExitCode != 0 {
		t.Fatalf("full multi-arg run exit=%d\nstderr:\n%s\nstdout:\n%s",
			resp.FullExitCode, resp.FullStderr, resp.FullStdout)
	}
	if resp.SubsetExitCode != 0 {
		t.Fatalf("subset run exit=%d\nstderr:\n%s\nstdout:\n%s",
			resp.SubsetExitCode, resp.SubsetStderr, resp.SubsetStdout)
	}

	manFull := resp.ManifestAfterFull
	if manFull == "" {
		t.Fatalf("expected %s after full run under %s", genManifestName, req.GenDir)
	}
	if !manifestHasTreePrefix(manFull, "tree-a") {
		t.Fatalf("full manifest missing tree-a entries:\n%s", manFull)
	}
	if !manifestHasTreePrefix(manFull, "tree-b") {
		t.Fatalf("full manifest missing tree-b entries:\n%s", manFull)
	}

	manSub := resp.ManifestAfterSubset
	if manSub == "" {
		t.Fatalf("expected %s after subset run under %s", genManifestName, req.GenDir)
	}
	// Core policy: subset must not shrink away the sibling tree's ledger entries.
	if !manifestHasTreePrefix(manSub, "tree-b") {
		t.Fatalf("subset run shrank manifest: tree-b entries missing (must keep sibling ledger)\n--- after full ---\n%s\n--- after subset ---\n%s",
			manFull, manSub)
	}
	// Sanity: selected tree still listed.
	if !manifestHasTreePrefix(manSub, "tree-a") {
		t.Fatalf("subset run lost tree-a entries:\n%s", manSub)
	}
	if !resp.SiblingGenDirExists {
		t.Fatalf("subset run removed sibling gen packages under %s/tree-b (must keep on disk)", req.GenDir)
	}
}
```
