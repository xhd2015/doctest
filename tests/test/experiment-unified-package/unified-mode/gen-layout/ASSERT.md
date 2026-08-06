## Expected

- Suite run succeeds (generation + go test of suite).
- Gen root contains `__droot`, `__registry`, `__allleaves`, and `suite`.
- At least one suite `*_test.go` exists.
- Fixture leaves `a`/`b` have **no** `*_test.go` under gen.
- At least two leaf non-test `.go` files exist under `a`/`b` paths.
- Each such leaf non-test file defines `func RunTestLeaf`.
- Suite **package** (all `suite/*.go`: `runall.go` + thin `suite_test.go`) imports
  only non-stdlib packages under `__registry` and `__allleaves` (plus stdlib).
  Fan-in may live in `runall.go`; `suite_test.go` may import only `testing`.

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("unified-mode RunTest failed before layout check: %s\nstderr:\n%s\ngen=%s",
			resp.RunErr, resp.Stderr, resp.GenDir)
	}
	if resp.GenDir == "" {
		t.Fatal("GenDir empty; cannot inspect unified layout")
	}
	if !resp.HasDroot {
		t.Fatalf("gen missing __droot (go files=%v gen=%s)", basenames(resp.GoFiles), resp.GenDir)
	}
	if !resp.HasRegistry {
		t.Fatalf("gen missing __registry (go files=%v gen=%s)", basenames(resp.GoFiles), resp.GenDir)
	}
	if !resp.HasAllLeaves {
		t.Fatalf("gen missing __allleaves (go files=%v gen=%s)", basenames(resp.GoFiles), resp.GenDir)
	}
	if !resp.HasSuite {
		t.Fatalf("gen missing suite package (go files=%v gen=%s)", basenames(resp.GoFiles), resp.GenDir)
	}
	if len(resp.SuiteTestFiles) < 1 {
		t.Fatalf("expected ≥1 suite *_test.go, got %v", basenames(resp.SuiteTestFiles))
	}
	if len(resp.LeafTestGoFiles) != 0 {
		t.Fatalf("unified leaves must not have *_test.go; got %v", basenames(resp.LeafTestGoFiles))
	}
	if len(resp.LeafNonTestGoFiles) < 2 {
		t.Fatalf("expected ≥2 leaf non-test .go under a/b, got %v (all go=%v)",
			basenames(resp.LeafNonTestGoFiles), basenames(resp.GoFiles))
	}
	for i, leaf := range resp.LeafNonTestGoFiles {
		if i < len(resp.LeafHasRunTestLeaf) && !resp.LeafHasRunTestLeaf[i] {
			t.Fatalf("leaf non-test %s missing func RunTestLeaf", filepath.ToSlash(leaf))
		}
	}
	if len(resp.SuiteNonTestFiles) < 1 {
		t.Fatalf("expected ≥1 suite non-test .go (runall.go), got test=%v non-test=%v",
			basenames(resp.SuiteTestFiles), basenames(resp.SuiteNonTestFiles))
	}
	if !suiteImportsOnlyRegistryAndAllLeaves(resp.SuiteImportLines) {
		t.Fatalf("suite package must import only __registry + __allleaves (plus stdlib); imports=%v suite test=%v non-test=%v",
			resp.SuiteImportLines, basenames(resp.SuiteTestFiles), basenames(resp.SuiteNonTestFiles))
	}
}

func basenames(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = filepath.Base(p)
	}
	return out
}
```
