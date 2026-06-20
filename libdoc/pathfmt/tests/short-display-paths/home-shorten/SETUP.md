# Scenario

**Feature**: paths outside cwd but under home display with `~` prefix

```
# home shorten (when not under cwd)
path under home -> "~" + suffix
```

## Preconditions

- Cwd is a temp directory outside the target home paths (typically under `os.TempDir()`).

## Steps

1. Save and restore cwd.
2. `chdir` to a fresh temp directory so target paths are not under cwd.

## Context

- Leaves set `req.Path` to home or a subdirectory of home.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	saveAndRestoreCwd(t)
	chdirTo(t, t.TempDir())
	return nil
}
```