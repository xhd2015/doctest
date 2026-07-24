# Scenario

**Feature**: multi-tree ./... with one prepare failure prints honest FAIL summary

```
# good tree prepares and passes; bad tree fails prepare (invalid Go in DOCTEST)
doctest test ./... --gen-dir <tmp> --no-color
  -> survivors run -> FAIL(p/t) on stdout (never PASS when prepare failed)
  -> process error includes "prepare failed:"
```

## Preconditions

- One module root with sibling trees `good/` (1 pass leaf) and `bad/` (syntax error in DOCTEST.md).
- Workdir is the module root so `./...` discovers both trees.

## Steps

1. Create prepare-fail multi-tree module fixture.
2. Run `doctest test --gen-dir <tmp> --no-color ./...` from module root.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	mod := createPrepareFailMultiTree(t)
	genDir := filepath.Join(t.TempDir(), "gen")
	// WorkDir + ./... needs subprocess isolation (L2 ignores WorkDir).
	req.UseCLI = true
	req.WorkDir = mod
	req.Args = []string{"test", "--gen-dir", genDir, "--no-color", "./..."}
	return nil
}
```
