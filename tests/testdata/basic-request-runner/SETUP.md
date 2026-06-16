# Scenario

**Feature**: tests for basic request runner

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Global Setup

## Context

- The root defines the shared `Request` and `Response` model.
- The root `Run` handles the default operation for all descendants.

```go
import "fmt"

type Request struct {
	Action string
	Name   string
}

type Response struct {
	Message string
}

func Setup(t *testing.T, req *Request) error {
	req.Name = "world"
	return nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	switch req.Action {
	case "greet":
		return &Response{Message: fmt.Sprintf("hello %s", req.Name)}, nil
	case "fail":
		return nil, fmt.Errorf("requested failure for %s", req.Name)
	default:
		return nil, fmt.Errorf("unknown action %q", req.Action)
	}
}
```
