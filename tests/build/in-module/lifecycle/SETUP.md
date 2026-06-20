# Scenario

**Feature**: compile temp lifecycle — created for internal imports, removed after run

```
# temp compile dir lifecycle
.doctest_run_* created under moduleRoot -> go test -> temp removed; dump kept
```

## Preconditions

- Tests verify `.doctest_run_*` directories are cleaned up after doctest completes.

## Steps

1. Prefix doctest args with the `test` subcommand.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"test"}, req.Args...)
	req.Env = append(req.Env, "GOWORK=off")
	return nil
}
```