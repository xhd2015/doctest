# Scenario

**Feature**: the single-agent dev-test skill is available through the skill registry

```
doctest skill dev-test --show -> embedded self-contained workflow -> stdout
```

## Preconditions

- The `doctest-dev-test` document is embedded in the product.

## Steps

1. Run `doctest skill dev-test --show`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "dev-test", "--show"}
	return nil
}
```
