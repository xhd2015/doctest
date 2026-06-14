## Preconditions
- The doctest binary is built by root Setup.
- Timeout is set to 120s.

## Steps
1. Prefix args with "test" command.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```
