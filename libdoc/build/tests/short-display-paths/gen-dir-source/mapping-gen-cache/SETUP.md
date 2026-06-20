# Scenario

**Feature**: auto mapping-gen cache dir displays with `~/...mapping-gen...`

```
# gen-dir modes
auto gen-dir -> mapping-gen cache under home

# stderr call sites
announceRoots -> DisplayPath(genRoot) | cd preview -> DisplayPath(runDir)
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