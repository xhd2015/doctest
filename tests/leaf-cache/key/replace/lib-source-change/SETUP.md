# Scenario

**Feature**: changing replace-target library source changes the key

```
# ../lib/lib.go
Val()=7  -> key1
Val()=100 -> key2
# key1 != key2
```

## Preconditions

- Replace fixture; ASSERT imports `example.com/lib`.
- Mutation = `replace_lib_src`.

## Steps

1. Set Mutation to `replace_lib_src`.
2. Run compute_mutate.
3. Assert keys differ.

## Context

- Source under a local replace path is first-class local content.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	req.Mutation = "replace_lib_src"
	return nil
}
```
