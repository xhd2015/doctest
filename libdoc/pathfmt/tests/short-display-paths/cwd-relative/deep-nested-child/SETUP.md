# Scenario

**Feature**: deeply nested child displays as `"./a/b/c"`

```
# cwd rules (checked first)
strict child of cwd -> "./" + rel
```

## Steps

1. Set `req.Path` to `a/b/c` under project root.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Path = filepath.Join(req.Path, "a", "b", "c")
	return nil
}
```