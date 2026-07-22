# Scenario

**Bug**: closed stdin file causes Stat error during agent implement

```
closed stdin file inject -> readStdinIfPresent Stat -> error returned
```

## Preconditions
- A closed file is injected as stdin, causing `Stat()` to fail.

## Steps
1. Open `/dev/null`, close it, assign to `req.StdinFile`.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	f, err := os.Open(os.DevNull)
	if err != nil {
		return err
	}
	f.Close()
	req.StdinFile = f
	return nil
}
```
