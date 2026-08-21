# Scenario

**Feature**: version and show actions are mutually exclusive

```
doctest skill dev-test --version --show -> action conflict
```

## Steps

1. Request both version and show.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "dev-test", "--version", "--show"}
	return nil
}
```
