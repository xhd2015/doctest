# Scenario

**Feature**: embed script emits single file with package assert and no test sources

```
# one output file
script/embed-assert -> assert.go containing package assert only
```

## Steps

1. Run embed script against `assert/` directory.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.SecondRun = false
	return nil
}
```