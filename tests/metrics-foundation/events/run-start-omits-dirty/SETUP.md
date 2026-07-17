# Scenario

**Feature**: run_start may include branch/commit but never a dirty-worktree field

```
# git metadata allowed
run_start{git_branch, git_commit}  # no dirty / git_dirty / dirty_worktree
```

## Preconditions

- Single run_start event with branch and commit set.

## Steps

1. Write one run_start with git_branch and git_commit.
2. Decode and assert dirty-related keys are absent.

## Context

- P1 explicitly excludes dirty git field from the schema.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Events = []map[string]any{
		{
			"type":       "run_start",
			"run_id":     "run-nodirty-1",
			"ts":         "2026-07-17T08:00:00Z",
			"project_id": "github.com_xhd2015_doctest",
			"cwd":        "/repo",
			"argv":       []any{"doctest", "test"},
			"git_branch": "feature/metrics",
			"git_commit": "deadbeefcafebabe",
		},
	}
	return nil
}
```
