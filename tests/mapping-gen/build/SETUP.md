## Preconditions
- The doctest binary is built by root Setup.
- Timeout is set to 120s.

## Steps
1. Prefix args with "build" command.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"build"}, req.Args...)
	return nil
}
```
