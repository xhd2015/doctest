# Scenario

**Feature**: the vet command detects anti-patterns in test file content

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | DOCTEST_SESSION_ID env read | assert without setup | skipped testdata
```

## Preconditions
- The vet command detects anti-patterns in test file content.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
