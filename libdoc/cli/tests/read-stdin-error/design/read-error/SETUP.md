# Scenario

**Bug**: directory stdin causes ReadAll error during agent design

```
directory as stdin inject -> readStdinIfPresent ReadAll -> error returned
```

## Preconditions
- An open directory is injected as stdin, causing `io.ReadAll()` to fail while `Stat()` succeeds.

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
