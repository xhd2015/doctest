# Scenario

**Feature**: fixed WARNING text for slow default suite

```
# pure formatter (no I/O)
FormatDefaultSuiteSlowWarning() -> string starting with WARNING:
```

## Preconditions

- Message content is stable for docs and skill cross-links.

## Steps

1. Call FormatDefaultSuiteSlowWarning.
2. Assert required phrases.

## Context

- Intent (may match closely):
  `WARNING: doctest default suite should be fast (run within 3 minutes); you're strongly recommended to use skill:doctest-review-perf to optimize default test suite performance (doctest skill review-perf --show)`

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "format_warn"
	return nil
}
```
