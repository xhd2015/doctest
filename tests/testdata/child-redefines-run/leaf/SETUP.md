# Scenario

**Feature**: override the Action field

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

## Steps
1. Override the Action field.
2. Define Run again — this redefines Run and is the violation under test.

```go
import (
	"fmt"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Action = "leaf-action"
	return nil
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	return &Response{Result: fmt.Sprintf("child:%s", req.Action)}, nil
}
```