# Scenario

**Feature**: `--no-metrics` sets NoMetrics opt-out

```
# parse with flag
parseTestOptions(["--no-metrics", "./tests"]) -> NoMetrics=true
```

## Preconditions

- Flag may appear before the directory operand.

## Steps

1. Parse `--no-metrics` plus a path.
2. Expect NoMetrics true; remain args still include the path.

## Context

- Recording integration leaf separately proves no files are written.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--no-metrics", "./tests"}
	return nil
}
```
