# Scenario

**Feature**: unknown runner flags should fail at parse time (L2)

```
runner.ParseTestOptions([dir, --definitely-not-real])
  -> error: unrecognized flag
```

## Preconditions

- Nested L2 root: `runner.ParseTestOptions` only; no product binary.
- Directory operand is unused once parse fails.

## Steps

1. Set `req.Args` to a dummy dir plus an unknown flag.
2. Run parses; Assert checks non-zero exit style and message.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Dir is unused: parse rejects the unknown flag first.
	req.Args = []string{".", "--definitely-not-real"}
	return nil
}
```
