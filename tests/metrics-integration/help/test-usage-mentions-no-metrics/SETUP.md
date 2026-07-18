# Scenario

**Feature**: `doctest test --help` documents `--metrics-on`

```
doctest test --help -> Options include --metrics-on (opt in to suite metrics JSONL)
```

## Preconditions

- Flag is already parsed by `ParseTestOptions` (P2); help text must list it.

## Steps

1. Run `doctest test --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"test", "--help"}
	return nil
}
```
