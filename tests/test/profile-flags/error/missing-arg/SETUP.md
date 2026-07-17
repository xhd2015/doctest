# Scenario

**Feature**: `-cpuprofile` without a value fails at CLI parse time

```
doctest test -cpuprofile
  -> parse error -> non-zero exit -> stderr mentions cpuprofile
```

## Preconditions
- No profile path argument follows `-cpuprofile`.

## Steps
1. Run `doctest test -cpuprofile` (no value, no dir required for parse failure).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Missing value for -cpuprofile; parse should fail before dir handling.
	req.Args = []string{"test", "-cpuprofile"}
	return nil
}
```
