## Steps
- Create an empty temp directory with no doctest marker files.
- Set `req.Input` to that directory.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Input = t.TempDir()
	return nil
}
```
