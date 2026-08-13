# Scenario

**Feature**: leaf calls the goto helper so it is not unused

```
leaf Setup -> installFakeOpencode
```

## Steps

1. Call `installFakeOpencode` so the helper is referenced.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := t.TempDir()
	if err := installFakeOpencode(filepath.Join(dir, "fake-opencode")); err != nil {
		return err
	}
	_ = os.Stdout
	return nil
}
```
