# Scenario

**Feature**: `doctest skill designer show` includes the resolved spec version

```
# designer prompt served to stdout
doctest skill designer show -> PROMPT.md with version 0.0.2
```

## Preconditions

- The designer prompt embed references `__DOCTEST_VERSION__` in its template.

## Steps

1. Run `doctest skill designer show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "designer", "show"}
	return nil
}
```