# Scenario

**Feature**: on the same Replace key, pre_test hook overwrites the user seed

```
user seed: project-source -> "from-user"
           active-vendor  -> "seed-only"   # not touched by hook
hook write: project-source -> "from-hook"
  -> final[project-source] == "from-hook"  (hook wins)
  -> final[active-vendor]  == "seed-only"  (proves seed was applied)
```

## Preconditions

- User seed is applied **before** the hook runs.
- A second seed-only key proves the seed layer ran (avoids false-green if only
  the hook wrote and the contested key happened to match).

## Steps

1. Seed two alias keys; hook overwrites only `project-source`.
2. Run expands aliases to fixture absolute paths.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Alias keys expanded by Run to fixture absolute paths before writing user JSON.
	req.UserReplace = map[string]string{
		"project-source": "from-user",
		"active-vendor":  "seed-only",
	}
	req.HookOverlays = [][]OverlayEntry{
		{{Source: "project-source", Replace: "from-hook"}},
	}
	return nil
}
```
