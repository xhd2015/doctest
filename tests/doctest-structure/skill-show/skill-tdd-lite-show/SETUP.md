# Scenario

**Feature**: `doctest skill tdd-lite --show` includes the resolved spec version

```
# TDD lite skill document served to stdout
doctest skill tdd-lite --show -> DOCTEST_TDD_LITE.md with version 0.0.2
```

## Preconditions

- The TDD lite skill embed references `__DOCTEST_VERSION__` in its template.

## Steps

1. Run `doctest skill tdd-lite --show`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "tdd-lite", "--show"}
	return nil
}
```