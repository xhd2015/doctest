# Scenario

**Feature**: this fixture tree has a root SETUP.md that defines Request, Response, and Setup,

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Root Missing Run Fixture

## Preconditions
- This fixture tree has a root SETUP.md that defines Request, Response, and Setup,
  but deliberately has no Run function.
- The missing Run should trigger a validation error during tree discovery.

## Steps
1. Build the tree. Discovery should fail because the root has no Run.

```go
import "fmt"

func Setup(t *testing.T, req *Request) error {
	if req.Input == "" {
		return fmt.Errorf("Input is required")
	}
	return nil
}
```
