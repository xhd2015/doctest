# Scenario

**Feature**: remote module sources are not file-hashed

```
# go.mod may require example.com/remote v1.0.0 (identity only)
# files under remote-src/... are NOT on a local replace path

remote-src/.../remote.go change -> key stable
```

## Preconditions

- Flavor = `remote`.
- A fake tree under `WorkDir/remote-src/...` looks like module cache content but
  is not linked via replace and is not under ModuleRoot.

## Steps

1. Rebuild workspace with remote flavor.
2. Mutate the remote-like file; key must stay stable.

## Context

- Only go.mod/go.sum identity may reflect remote modules — never their source trees.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Flavor = "remote"
	req.Op = "compute_mutate"
	req.WorkDir = t.TempDir()
	req.ModuleRoot = ""
	req.TreeRoot = ""
	req.LeafDir = ""
	return ensureWorkspace(t, req)
}
```
