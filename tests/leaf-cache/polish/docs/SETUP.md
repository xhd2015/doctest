# Scenario

**Feature**: CLI help documents leaf-cache control flags

```
doctest test --help
  -> mentions -a
  -> mentions --no-leaf-cache
```

## Preconditions

- Binary from polish root Setup (testbin.Ensure).

## Steps

1. Op=`runtime_once` with Args `test --help` (or `test`, `-h`).
2. Assert stdout/stderr contain flag names.

## Context

- Documentation contract so users discover disable knobs.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_once"
	req.Args = []string{"test", "--help"}
	// Help should not need leaf-cache env isolation, but Env is already set.
	return nil
}
```
