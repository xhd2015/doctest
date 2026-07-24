# Scenario

**Feature**: `--metrics-on` sets MetricsOn opt-in

```
# parse with flag
parseTestOptions(["--metrics-on", "./tests"]) -> MetricsOn=true
```

## Preconditions

- Flag may appear before the directory operand.

## Steps

1. Parse `--metrics-on` plus a path.
2. Expect MetricsOn true; remain args still include the path.

## Context

- Recording integration leaf separately proves no files are written.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"--metrics-on", "./tests"}
	return nil
}
```
