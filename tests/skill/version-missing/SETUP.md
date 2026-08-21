# Scenario

**Feature**: querying an unversioned skill fails explicitly

```
doctest skill tdd --version -> missing metadata.version error
```

## Steps

1. Query the version of the existing unversioned `tdd` skill.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"skill", "tdd", "--version"}
	return nil
}
```
