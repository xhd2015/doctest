# Scenario

**Feature**: explicit `--gen-dir` under sandbox project without process Chdir

```
# gen-dir modes
explicit gen-dir -> req.GenDir="_gen" joined to abs projRoot/_gen in Run
  # NOT: os.Chdir(projRoot); GenDir="_gen" relative to process cwd

# stderr call sites
announceRoots -> pathfmt.Short(absGen)
cd preview -> pathfmt.Short(runDir under absGen)
doctest: -> pathfmt.Short(testRoot)
```

## Steps

1. Set `req.GenDir` to `"_gen"` (relative to sandbox project root; Run makes absolute).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.GenDir = "_gen"
	return nil
}
```
