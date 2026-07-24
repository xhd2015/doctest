# Scenario

**Feature**: set the request Name to "doctest"

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

## Steps
1. Set the request Name to "doctest".

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Name = "doctest"
    return nil
}
```
