# Scenario

**Feature**: `doctest test --help` lists `--label`

## Steps

1. Run help subcommand.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"test", "--help"}
	return nil
}
```