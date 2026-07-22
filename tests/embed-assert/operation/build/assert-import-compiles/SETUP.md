# Scenario

**Feature**: doctest build succeeds when leaf imports assert package

```
# doctest build
doctest build <tests> -> generated code compiles with assert replace
```

## Preconditions

- Public module with assert-importing leaf.

## Steps

1. Create public module with assert leaf.
2. Run `doctest build <tests>`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t, req, "", defaultAssertAssertGo())
	setupModuleEnv(t, req)
	req.Args = []string{"build", req.TestDir}
	return nil
}
```