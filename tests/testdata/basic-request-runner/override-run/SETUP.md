# Scenario

**Feature**: mutate the inherited request to use the root Run

```
# tree structure validation during build/test
root: must define Request, Response, Run
child: must define Setup, must NOT redefine Run
leaf: ASSERT.md with func Assert
```

# Override Run Setup

## Steps

1. Mutate the inherited request to use the root Run.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Action = "greet"
	req.Name = "leaf"
	return nil
}
```
