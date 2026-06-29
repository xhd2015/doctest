# Scenario

**Feature**: `doctest skills --help` lists update subcommand

```
doctest skills --help -> stdout contains "update"
```

## Preconditions

- None.

## Steps

1. Run with `["skills", "--help"]`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skills", "--help"}
	return nil
}
```