# Scenario

**Feature**: two root helpers where the caller is defined before the callee

```
# Func-literal assembly during code generation
writeHelperDecls -> topo sort -> generated test compiles
```

## Preconditions

- Root `DOCTEST.md` defines two helpers in source order:
  `func caller(x int) string { return callee(x) }` followed by
  `func callee(x int) string { return "x" }`.
- This ordering is legal for top-level funcs (forward references allowed) but
  illegal for func literals: when `writeHelperDecls` emits helpers verbatim in
  source order, `caller := func(x int) string { return callee(x) }` references
  `callee` before it is declared, producing `undefined: callee`.
- A leaf `Setup` calls `caller` so doctest must lower both helpers to func
  literals.

## Steps

1. Create a temp Go project with `go.mod`.
2. Create a doctest root whose `DOCTEST.md` go block includes `caller` then
   `callee` (caller first).
3. Create a leaf whose `Setup` calls `caller(42)` and checks the result.
4. Run `doctest test -v <test-dir>` with the module's doctest binary.
5. ASSERT expects generation and execution to succeed (exit code 0): the codegen
   must topologically order helper closures so callees precede callers.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	testDir := filepath.Join(proj, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}

	runCode := `func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) { return &Response{}, nil }

func caller(x int) string { return callee(x) }
func callee(x int) string { return "x" }
`
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}

	leafDir := filepath.Join(testDir, "leaf")
	leafSetup := doctestGoBlock(`import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if got := caller(42); got != "x" {
		t.Fatalf("caller: got %q", got)
	}
	_ = req
	return nil
}`)
	if err := createDoctestLeaf(leafDir, leafSetup); err != nil {
		t.Fatalf("create leaf: %v", err)
	}

	req.Args = []string{"test", "-v", testDir}
	return nil
}
```
