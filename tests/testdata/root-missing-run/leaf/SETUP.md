# Scenario

**Feature**: set the Input field on the request

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

## Steps
1. Set the Input field on the request.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Input = "hello"
    return nil
}
```
