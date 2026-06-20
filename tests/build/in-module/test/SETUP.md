# Scenario

**Feature**: doctest test command with internal import scan and temp compile

```
# internal import detected -> temp compile under moduleRoot
doctest test <tree> -> scan imports -> .doctest_run_* -> go test
```

## Preconditions

- Operation mode is `doctest test`.

## Steps

1. Prefix doctest args with the `test` subcommand.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```