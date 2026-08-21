# Scenario

**Feature**: version cannot be combined with header mode

```
doctest skill dev-test --version --header -> header/show constraint error
```

## Steps

1. Request version with the show-only header modifier.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "dev-test", "--version", "--header"}
	return nil
}
```
