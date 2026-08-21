# Scenario

**Feature**: version lookup accepts name-first arguments

```
doctest skill dev-test --version -> 0.1.0
```

## Steps

1. Run `doctest skill dev-test --version`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "dev-test", "--version"}
	return nil
}
```
