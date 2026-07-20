# Scenario

**Feature**: changing a remote-looking source file outside the local DAG does not change the key

```
# before
ComputeLeafKey -> key1

# edit WorkDir/remote-src/example.com/remote@v1.0.0/remote.go
ComputeLeafKey -> key2
# key1 == key2
```

## Preconditions

- Remote flavor fixture from parent.
- Mutation = `remote_proxy_file`.

## Steps

1. Set Mutation to `remote_proxy_file`.
2. Run compute_mutate.
3. Assert Key == Key2.

## Context

- If this key flips, the implementation is walking non-local trees (bug).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "compute_mutate"
	req.Mutation = "remote_proxy_file"
	return nil
}
```
