# Scenario

**Feature**: a bare `...` pattern is used instead of `./...` or a qualified path

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- A bare `...` pattern is used instead of `./...` or a qualified path.

## Steps
1. Run `doctest vet ...`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"vet", "..."}
	return nil
}
```
