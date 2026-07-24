# Scenario

**Feature**: skills subcommand help documents update

```
doctest skills --help -> usage mentions update
```

## Preconditions

- Binary built and cwd set.

## Steps

1. Leaves set `req.Args` to skills help invocation.

## Context

- No filesystem side effects expected.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreInstalls = nil
	return nil
}
```