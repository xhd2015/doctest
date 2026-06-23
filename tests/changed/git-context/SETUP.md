# Scenario

**Feature**: `--changed` behavior depends on git repository presence

```
# resolve git root for change detection
doctest --changed -> git.ShowToplevel -> error if not in repo

# filter by on-disk changes
git.GetOnDiskChangedFiles -> map paths -> affected leaves
```

## Preconditions

- The doctest binary is built.

## Steps

1. Choose whether the working directory is inside a git repository.
2. Configure fixture layout and subcommand args in descendant setups.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "CHANGED_TEST_GROUP=git-context")
	return nil
}
```