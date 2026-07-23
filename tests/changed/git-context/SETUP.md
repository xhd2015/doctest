# Scenario

**Feature**: `--changed` behavior depends on git-context vs pure policy

```
# L2 policy (in-git-repo/*): synthetic changed paths, no binary
fixture tree -> core.FilterByChangedFiles | ChangedDoctestMarkdownFiles

# L3 process (not-git-repo/*): real CLI requires git repo
doctest --changed outside .git -> hard error
```

## Preconditions

- Descendants choose L2 policy fixtures or L3 CLI.

## Steps

1. Choose whether the scenario is pure selection policy (L2) or process-boundary (L3).
2. Configure fixture layout and request fields in descendant setups.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	_ = req
	return nil
}
```
