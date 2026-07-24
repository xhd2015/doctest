# Scenario

**Feature**: changing replace-target go.mod changes the key

```
# ../lib/go.mod
go 1.22 -> key1
go 1.23 -> key2
# key1 != key2
```

## Preconditions

- Replace fixture.
- Mutation = `replace_lib_gomod`.

## Steps

1. Set Mutation to `replace_lib_gomod`.
2. Run compute_mutate.
3. Assert keys differ.

## Context

- Requirement: local replace modules' go.mod is context for the DAG, even when
  `.go` sources are unchanged.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	req.Mutation = "replace_lib_gomod"
	return nil
}
```
