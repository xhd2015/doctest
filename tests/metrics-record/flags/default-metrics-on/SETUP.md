# Scenario

**Feature**: metrics recording enabled by default (no opt-out flag)

```
# parse without --no-metrics
parseTestOptions(["./tests"]) -> NoMetrics=false
```

## Preconditions

- Remain args may include the directory operand.

## Steps

1. Parse args with only a directory path.
2. Expect NoMetrics false.

## Context

- Default: metrics **on**.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"./tests"}
	return nil
}
```
