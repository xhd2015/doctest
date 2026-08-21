# Scenario

**Feature**: version lookup accepts action-first arguments

```
doctest skill --version dev-test -> 0.1.0
```

## Steps

1. Run `doctest skill --version dev-test`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "--version", "dev-test"}
	return nil
}
```
