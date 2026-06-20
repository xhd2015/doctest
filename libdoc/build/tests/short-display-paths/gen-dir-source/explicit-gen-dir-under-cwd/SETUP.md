# Scenario

**Feature**: explicit `--gen-dir` under project cwd displays as `→ ./_gen`

```
# gen-dir modes
explicit gen-dir -> user path under cwd

# stderr call sites
announceRoots -> DisplayPath(genRoot)
```

## Steps

1. Set `req.GenDir` to `"_gen"` (relative to project root / cwd).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.GenDir = "_gen"
	return nil
}
```