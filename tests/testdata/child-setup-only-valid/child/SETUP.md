# Scenario

**Feature**: the root provides Request, Response, and Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

## Preconditions
- The root provides Request, Response, and Run.
- This child defines Setup only, which is the expected pattern.

## Steps
1. Set the Value on the request.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Value = 7
    return nil
}
```
