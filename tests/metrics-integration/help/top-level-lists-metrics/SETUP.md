# Scenario

**Feature**: top-level help lists the `metrics` command

```
doctest --help -> Usage … Metrics: metrics path|last|top|…
```

## Preconditions

- None beyond binary.

## Steps

1. Run `doctest --help`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"--help"}
	return nil
}
```
