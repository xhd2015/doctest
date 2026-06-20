# Scenario

**Feature**: strict child of cwd displays as `"./child"`

```
# cwd rules (checked first)
strict child of cwd -> "./" + rel
```

## Steps

1. Set `req.Path` to the `child` subdirectory under project root.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Path = filepath.Join(req.Path, "child")
	return nil
}
```