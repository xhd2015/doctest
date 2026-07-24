# Scenario

**Feature**: no additional changes to the request; Value already set by parent

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

## Steps
1. No additional changes to the request; Value already set by parent.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    _ = req
    return nil
}
```
