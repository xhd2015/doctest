# Scenario

**Feature**: path equal to cwd displays as `"."`

```
# cwd rules (checked first)
path == cwd -> "."
```

## Steps

1. Set `req.Path` to the current working directory (project root).

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	req.Path = wd
	return nil
}
```