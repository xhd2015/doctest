# Scenario

**Feature**: show-header preserves standard nested version metadata

```
doctest skill dev-test --show --header -> YAML metadata.version
```

## Steps

1. Print only the `dev-test` YAML frontmatter.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "dev-test", "--show", "--header"}
	return nil
}
```
