# Scenario

**Feature**: a minimal valid doctest tree is available

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- A minimal valid doctest tree is available.

## Steps
1. Run `doctest vet -v <dir>`.

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
	req.Args = []string{"vet", "-v", dir}
	return nil
}
```
