## Preconditions
- testdata/ has go.mod at root with a DOCTest tree at group/subgroup/.

## Steps
1. Set WorkDir to the testdata/ directory.
2. Run `doctest test -v ./group/subgroup/...`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.WorkDir = "./testdata"
    req.Args = []string{"test", "-v", "./group/subgroup/..."}
    return nil
}
```
