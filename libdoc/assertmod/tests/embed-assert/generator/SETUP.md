# Scenario

**Feature**: embed script concatenates assert sources into one generated file

```
# script/embed-assert
assert/*.go (sorted, no *_test.go) -> single assert.go bytes
```

## Preconditions

- Siblings test output shape (A1) and determinism (A2).

## Steps

1. Descendant sets `runKind = "embed-script"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	runKind = "embed-script"
	req.OutputPath = ""
	return nil
}
```