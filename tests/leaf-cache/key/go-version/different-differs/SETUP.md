# Scenario

**Feature**: different GoVersion values produce different keys for the same spine

```
in + go1.25.0 -> key1
in + go1.24.0 -> key2
# key1 != key2
```

## Preconditions

- Parent set Op and both version strings.
- Fixture sources are unchanged between the two calls.

## Steps

1. Confirm GoVersion / GoVersionB remain distinct.
2. Run compute_go_versions.
3. Assert keys differ.

## Context

- Go version is context for the DAG, not part of the file tree.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_go_versions"
	if req.GoVersion == req.GoVersionB {
		req.GoVersion = "go1.25.0"
		req.GoVersionB = "go1.24.0"
	}
	return nil
}
```
