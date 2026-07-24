# Scenario

**Feature**: changing an ancestor SETUP Go changes the key

```
# original group SETUP Go
group/SETUP.md Go -> key1

# mutate group Setup body
group/SETUP.md Go' -> key2
# key1 != key2
```

## Preconditions

- Base fixture includes intermediate `group/SETUP.md` between root and leaf.
- Mutation = `ancestor_setup`.

## Steps

1. Set Mutation to `ancestor_setup`.
2. Run compute_mutate.
3. Assert Key ≠ Key2.

## Context

- Ancestor setup is shared spine; editing it must invalidate every descendant leaf key.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	req.Mutation = "ancestor_setup"
	return nil
}
```
