# Scenario

**Feature**: warm second generate increases unchanged count in result summary

```
run1 cold -> modified high
run2 same GenDir identical content
  -> gen-plan: result summary unchanged increases (or modified decreases)
```

## Preconditions

- Same fixture and GenDir for two CLI invocations (`Mode=cli-twice`).

## Steps

1. Switch Mode to cli-twice; Args2 mirrors Args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Mode = "cli-twice"
	req.Args2 = append([]string(nil), req.Args...)
	return nil
}
```
