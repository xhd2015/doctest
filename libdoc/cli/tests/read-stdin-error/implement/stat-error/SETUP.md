## Preconditions
- `os.Stdin` is replaced with a closed file, causing `os.Stdin.Stat()` to fail.

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
