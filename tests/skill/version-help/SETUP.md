# Scenario

**Feature**: help takes precedence over a version request

```
doctest skill --version --help -> skill-level help
```

## Steps

1. Request help together with version.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "--version", "--help"}
	return nil
}
```
