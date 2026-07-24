# Scenario

**Feature**: classic sets `DOCTEST_CASE` to the absolute leaf case directory

```
# TreeCase.Path = "nested/leaf"
# DocTestRoot = <abs root>
d.DOCTEST_CASE == filepath.Join(DocTestRoot, "nested/leaf")
```

## Preconditions

- Fixed abs root via `req.DocTestRoot` so assertions can match the literal path.
- Case path is `nested/leaf`.

## Steps

1. Assemble classic with known root and leaf path.
2. Assert generated source embeds that absolute CASE path under `DOCTEST_CASE`.

## Context

- When Path is empty, CASE would equal ROOT; this leaf covers the non-empty join form.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	abs, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	req.DocTestRoot = abs
	req.CasePath = "nested/leaf"
	req.AuthorDMode = "named-d"
	return nil
}
```
