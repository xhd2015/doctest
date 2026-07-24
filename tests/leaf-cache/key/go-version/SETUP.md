# Scenario

**Feature**: GoVersion string is mixed into the leaf key

```
# same module/tree/leaf
GoVersion=go1.25.0 -> keyA
GoVersion=go1.24.0 -> keyB
# keyA != keyB
```

## Preconditions

- Base fixture; no source mutations.
- Op is `compute_go_versions`.

## Steps

1. Set two distinct GoVersion strings.
2. Compute keys for each; they must differ.

## Context

- Different toolchains must not share pass markers even when sources match.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_go_versions"
	req.GoVersion = "go1.25.0"
	req.GoVersionB = "go1.24.0"
	return nil
}
```
