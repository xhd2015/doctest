# Scenario

**Feature**: doctest build compiles generated tests that import assert

```
# doctest build with assert import
doctest build <tree> -> go build succeeds with assert replace/modfile
```

## Preconditions

- Leaf ASSERT imports assert package.

## Steps

1. Descendant runs `doctest build` against assert-importing tree.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"build"}, req.Args...)
	return nil
}
```