# Scenario

**Feature**: select an action that the root `Run` intentionally reports as a run error

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Expected Error Setup

## Steps

1. Select an action that the root `Run` intentionally reports as a run error.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Action = "fail"
	req.Name = "case"
	return nil
}
```
