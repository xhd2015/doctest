# Scenario

**Feature**: root helper with named results whose body assigns to those names

```
# Func-literal assembly during code generation
writeFuncClosure -> imports.Process -> generated test compiles
```

## Preconditions

- Root `DOCTEST.md` defines a helper whose body assigns to its named result
  variables before returning, e.g.
  `func splitNames(req *Request) (mainRepo, wtDir, branch string) { mainRepo = "a"; ...; return }`.
- The current codegen lowers helpers to func literals using `ResultTypes`
  (names stripped), emitting
  `splitNames := func(req *Request) (string, string, string) { mainRepo = "a" ... }`.
- The body's assignments to `mainRepo`/`wtDir`/`branch` then reference undeclared
  identifiers, so the generated test fails to compile with `undefined: mainRepo`.
- A leaf `Setup` calls the helper so doctest must lower it to a func literal.

## Steps

1. Create a temp Go project with `go.mod`.
2. Create a doctest root whose `DOCTEST.md` go block includes `splitNames`.
3. Create a leaf whose `Setup` calls `splitNames(req)` and checks the values.
4. Run `doctest test -v <test-dir>` with the module's doctest binary.
5. ASSERT expects generation and execution to succeed (exit code 0): the closure
   signature must preserve named results so the body's assignments compile.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	proj := t.TempDir()
	if err := os.WriteFile(filepath.Join(proj, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	testDir := filepath.Join(proj, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir test dir: %v", err)
	}

	runCode := `func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }

func splitNames(req *Request) (mainRepo, wtDir, branch string) {
	mainRepo = "a"
	wtDir = "b"
	branch = "c"
	return
}`
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}

	leafDir := filepath.Join(testDir, "leaf")
	leafSetup := doctestGoBlock(`import "testing"

func Setup(t *testing.T, req *Request) error {
	mainRepo, wtDir, branch := splitNames(req)
	if mainRepo != "a" || wtDir != "b" || branch != "c" {
		t.Fatalf("splitNames: got %s,%s,%s", mainRepo, wtDir, branch)
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
