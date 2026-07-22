# Scenario

**Feature**: auto mapping-gen cache dir displays with `~/...mapping-gen...` without Chdir

```
# gen-dir modes
auto gen-dir -> mapping-gen cache under home

# stderr call sites (display-only Short; process cwd unchanged)
announceRoots -> pathfmt.Short(genRoot)
cd preview -> pathfmt.Short(runDir)   # ~/.../mapping-gen/...
doctest: -> pathfmt.Short(testRoot)   # abs temp when sandbox outside cwd
```

## Steps

1. Leave `req.GenDir` empty so `build.Test` picks the auto mapping-gen cache.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.GenDir = ""
	return nil
}
```
