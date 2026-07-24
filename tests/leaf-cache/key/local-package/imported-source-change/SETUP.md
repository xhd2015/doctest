# Scenario

**Feature**: editing a local package imported by the spine changes the key

```
# pkg/helper imported by ASSERT
helper.Answer body v1 -> key1
helper.Answer body v2 -> key2
# key1 != key2
```

## Preconditions

- Base fixture; leaf ASSERT imports `example.com/app/pkg/helper`.
- Mutation = `local_imported`.

## Steps

1. Set Mutation to `local_imported`.
2. Run compute_mutate.
3. Assert keys differ.

## Context

- Local import closure is the main reason keys track production code under test.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	req.Mutation = "local_imported"
	return nil
}
```
