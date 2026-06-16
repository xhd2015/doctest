# Scenario

**Feature**: select the default greeting action

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Happy Path Setup

## Steps

1. Select the default greeting action.
2. Override the request name for this leaf.

```go
func Setup(t *testing.T, req *Request) error {
	req.Action = "greet"
	req.Name = "runner"
	return nil
}
```
