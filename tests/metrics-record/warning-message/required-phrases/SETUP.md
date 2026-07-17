# Scenario

**Feature**: warning message includes WARNING banner, skill name, show command, and 3 minutes

```
# reader-facing content checklist
FormatDefaultSuiteSlowWarning()
  contains WARNING:
  contains skill:doctest-review-perf
  contains review-perf --show (or full doctest skill review-perf --show)
  contains 3 minutes
```

## Preconditions

- Single pure format call.

## Steps

1. Format the warning.
2. Assert substrings.

## Context

- Case-sensitive `WARNING:` prefix intent; skill identifier must appear for discoverability.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "format_warn"
	return nil
}
```
