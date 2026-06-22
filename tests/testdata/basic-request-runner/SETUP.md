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

func Setup(t *testing.T, req *Request) error {
	req.Name = "world"
	return nil
}
```
