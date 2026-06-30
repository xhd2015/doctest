# Scenario

**Feature**: root helper with multiple named results sharing one type `(port, alt int)`

```
# Func-literal assembly during code generation
writeFuncClosure -> imports.Process -> generated test compiles
```

## Preconditions

- Root `DOCTEST.md` defines `func pickTwoPorts(base int) (port, alt int)`.
- Go AST stores both names in a single result field, so `resultsString` omits
  outer parentheses and `writeFuncClosure` emits invalid syntax:
  `pickTwoPorts := func(base int) port int, alt int { ... }`.
- A leaf `Setup` calls the helper so doctest must lower it to a func literal.

## Steps

1. Create a temp Go project with `go.mod`.
2. Create a doctest root whose `DOCTEST.md` go block includes `pickTwoPorts`.
3. Create a leaf whose `Setup` calls `pickTwoPorts(req.Base)`.
4. Run `doctest test -v <test-dir>` with the module's doctest binary.
5. ASSERT expects generation and execution to succeed (exit code 0).

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

func pickTwoPorts(base int) (port, alt int) {
	return base, base + 1
}`
	if err := createDoctestRoot(testDir, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}

	leafDir := filepath.Join(testDir, "leaf")
	leafSetup := doctestGoBlock(`import "testing"

func Setup(t *testing.T, req *Request) error {
	port, alt := pickTwoPorts(7681)
	if port != 7681 || alt != 7682 {
		t.Fatalf("pickTwoPorts: got %d,%d", port, alt)
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