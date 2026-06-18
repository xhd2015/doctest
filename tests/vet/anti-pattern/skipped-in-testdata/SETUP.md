# Scenario

**Feature**: a doctest tree with an anti-pattern file inside a `testdata/` directory

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- A doctest tree with an anti-pattern file inside a `testdata/` directory.

## Steps
1. Create a valid doctest tree with a valid root SETUP.md.
2. Place a SETUP.md with embedded Go anti-pattern inside `testdata/`.
3. Run `doctest vet <dir>`.

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
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("# Scenario\n\n**Feature**: minimal test setup\n\n\x60\x60\x60\n# minimal pipeline\nsystem -> run\n\x60\x60\x60\n\n## Setup\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "testdata"), 0755); err != nil {
		t.Fatal(err)
	}
	fixture, err := os.ReadFile("fixture_setup.md.txt")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "testdata", "SETUP.md"), fixture, 0644); err != nil {
		t.Fatal(err)
	}
	req.Args = []string{"vet", dir}
	return nil
}
```
