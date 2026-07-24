# Scenario

**Feature**: `-cpuprofile` without a value fails at parse time (L2)

```
runner.ParseTestOptions([-cpuprofile])
  -> parse error mentioning cpuprofile / argument
```

## Preconditions

- Nested L2 root: parse only; no directory required for parse failure.

## Steps

1. Pass incomplete profile flag.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"-cpuprofile"}
	return nil
}
```
