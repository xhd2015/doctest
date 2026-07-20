# Scenario

**Feature**: content change refreshes go.mod bytes and manifest entry

```
drop replace localdep -> gen go.mod loses replace; manifest line for go.mod changes
```

## Steps

1. Inherit content-change Setup.
2. Confirm ChangeSourceGoMod is set and first go.mod had the replace.

```go
import (
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if req.Mode != "write-gomod-second" {
		t.Fatalf("updates-gomod-and-manifest expects Mode write-gomod-second, got %q", req.Mode)
	}
	if req.ChangeSourceGoMod == "" {
		t.Fatal("parent must set ChangeSourceGoMod for content-change path")
	}
	if !strings.Contains(snapGoModContentBefore, "replace localdep") {
		t.Fatalf("first gen go.mod should include replace localdep, got:\n%s", snapGoModContentBefore)
	}
	return nil
}
```
