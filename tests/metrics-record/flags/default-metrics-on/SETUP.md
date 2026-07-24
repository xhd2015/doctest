# Scenario

**Feature**: metrics are off by default (opt-in)

```
# parse without --metrics-on
parseTestOptions(["./tests"]) -> MetricsOn=false
```

## Preconditions

- No metrics-related flags in argv.

## Steps

1. Parse a directory-only argv.
2. Expect MetricsOn false.

## Context

- Recording is enabled only with `--metrics-on` / `MetricsOn=true`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"./tests"}
	return nil
}
```
