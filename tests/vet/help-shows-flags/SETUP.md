## Preconditions
- The vet command help should document the new `-v` flag and positional patterns.

## Steps
1. Run `doctest vet --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"vet", "--help"}
	return nil
}
```
