# Scenario

**Feature**: zero runnable leaves print no tests on stderr

```
# no leaves discovered
doctest test <empty-dir> -> Total=0 -> stderr "no tests"
```

## Preconditions

- Target directory has no doctest tree (empty temp dir).

## Steps

1. Create an empty temp directory.
2. Run `doctest test <empty-dir>`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	emptyDir := createEmptyDir(t)
	req.Args = []string{"test", emptyDir}
	return nil
}
```