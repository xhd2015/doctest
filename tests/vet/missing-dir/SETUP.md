# Scenario

**Feature**: no target directory argument is supplied

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- No target directory argument is supplied.

## Steps
1. Run `doctest vet`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"vet"}
    return nil
}
```
