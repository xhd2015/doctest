# Scenario

**Feature**: a doctest tree with a SETUP.md that embeds a Go program as a raw string literal

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- A doctest tree with a SETUP.md that embeds a Go program as a raw string literal.

## Steps
1. Create a minimal doctest tree with a SETUP.md containing a string literal that has `package main` and `func main()`.
2. Run `doctest vet <dir>`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# Tests\n\n## DSN (Domain Specific Notion)\n\n### Participants\n- **system** — under test.\n\n### Behaviors\n- **run** — executes the scenario.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	fixture, err := os.ReadFile("fixture_setup.md.txt")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), fixture, 0644); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"vet", dir}
	return nil
}
```
