# Scenario

**Feature**: `doctest test --help` lists `--label` (L2 CLI)

## Steps

1. Run help via in-process CLI.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"test", "--help"}
	return nil
}
```
