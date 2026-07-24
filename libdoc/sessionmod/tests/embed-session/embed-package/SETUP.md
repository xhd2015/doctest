# Scenario

**Feature**: sessionmod package accessors for embedded session sources

```
# accessor contract
Content() / ContentMD5() / RawSourceCacheKeyMD5()
```

## Steps

1. Leaf selects runKind.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Grouping node: ensure ModuleRoot already resolved by root Setup.
	if req.ModuleRoot == "" {
		t.Fatal("ModuleRoot must be set by root Setup")
	}
	return nil
}
```
