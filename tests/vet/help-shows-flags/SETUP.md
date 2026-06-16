# Scenario

**Feature**: the vet command help should document the new `-v` flag and positional patterns

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- The vet command help should document the new `-v` flag and positional patterns.

## Steps
1. Run `doctest vet --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"vet", "--help"}
	return nil
}
```
