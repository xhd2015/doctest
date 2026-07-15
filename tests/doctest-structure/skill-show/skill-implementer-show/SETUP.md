# Scenario

**Feature**: `doctest skill implementer --show` includes the resolved spec version

```
# implementer prompt served to stdout
doctest skill implementer --show -> PROMPT.md with version 0.0.2
```

## Preconditions

- The implementer prompt embed references `__DOCTEST_VERSION__` in its template.

## Steps

1. Run `doctest skill implementer --show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "implementer", "--show"}
	return nil
}
```