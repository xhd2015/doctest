# Scenario

**Bug**: directory stdin causes ReadAll error during agent implement

```
directory as os.Stdin -> readStdinIfPresent ReadAll -> error returned
```

## Preconditions
- `os.Stdin` is replaced with an open directory, causing `io.ReadAll()` to fail while `Stat()` succeeds.

## Steps
1. Open a temp directory, assign the file descriptor to `req.StdinFile`.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	f, err := os.Open(t.TempDir())
	if err != nil {
		return err
	}
	req.StdinFile = f
	return nil
}
```
