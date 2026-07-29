# Scenario

**Feature**: `--go-cmd=foo` is rejected (not treated as auto)

```
ParseTestOptions(["--go-cmd=foo", "."])
  -> error (invalid value / must be auto|xgo|go)
```

## Preconditions

- Value `foo` is not a valid mode.
- Error must describe invalid **value** (or list allowed values), not only a
  generic parse failure without naming go-cmd.

## Steps

1. Set parse args with `--go-cmd=foo` and a dummy directory remain.
2. Parse only.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ParseOnly = true
	// Equals form; implementer may also accept --go-cmd foo.
	req.ParseArgs = []string{"--go-cmd=foo", "."}
	return nil
}
```
