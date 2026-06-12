## Preconditions
- testdata/ has go.mod at root with two DOCTest trees: alpha/ and beta/.

## Steps
1. Set WorkDir to the testdata/ directory.
2. Run `doctest test -v alpha/...` (without `./` prefix).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.WorkDir = "./testdata"
    req.Args = []string{"test", "-v", "alpha/..."}
    return nil
}
```
