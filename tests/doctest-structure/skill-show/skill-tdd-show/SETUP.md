# Scenario

**Feature**: `doctest skill tdd --show` includes the resolved spec version

```
# TDD skill document served to stdout
doctest skill tdd --show -> DOCTEST_TDD.md with version 0.0.2
```

## Preconditions

- The TDD skill embed references `__DOCTEST_VERSION__` in its template.

## Steps

1. Run `doctest skill tdd --show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "tdd", "--show"}
	return nil
}
```