# Scenario

**Feature**: help surfaces mention metrics product commands and flags

```
doctest --help -> lists metrics
doctest test --help -> documents --no-metrics
```

## Preconditions

- Binary built via root Setup.
- No MetricsRoot fixtures required.

## Steps

1. Leaf sets `Op=help` and `Args`.
2. Run CLI help.
3. Assert substrings on stdout.

## Context

- Top-level usage currently documents a Metrics section; test usage should
  document `--no-metrics` (may be RED until help text is updated).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "help"
	return nil
}
```
