# Scenario

**Feature**: cached re-run grays the Cached summary segment

```
# run generated go tests, emit progress
build.Test(dir, opts) -> go test -> dot per package -> summary line

# summary coloring
summary Cached -> gray always (even when N>0)
```

## Preconditions
- One passing package; cache warmed by a prior `build.Test` on the same gen dir.

## Steps
1. Enable `WarmCache` and use one passing leaf.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.PassCount = 1
	req.FailCount = 0
	req.WarmCache = true
	return nil
}
```