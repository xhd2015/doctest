# Scenario

**Feature**: createDoctestTree helper writes a valid tree (**L2 helper**)

```
createDoctestTree(dir, stub=false)
  -> DOCTEST.md + SETUP.md + basic/{SETUP,ASSERT}.md on disk
```

## Preconditions

- This leaf only exercises the fixture helper used by e2e orchestrator leaves.
- No product binary / agent run required.

## Steps

1. Call createDoctestTree to write the test tree to a temp dir.
2. Stash path on req.SessionHome for Assert (parent-side; not process Env).
3. Run a no-op in-process short path so the harness completes.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := filepath.Join(t.TempDir(), "test-tree")
	createDoctestTree(t, dir, false)
	// Parent-side path for Assert — not process Env (UseCLI=false).
	req.SessionHome = dir
	req.UseCLI = false
	req.Bin = ""
	req.Env = nil
	req.Args = []string{"help"}
	return nil
}
```
