# Scenario

**Feature**: ref leaf injects `d *session.Doctest` and root package has no free DOCTEST_* vars

```
# AssembleRefLeafTestSource
leaf *_test.go:
  d := &session.Doctest{...}
  RootSetup / leaf setup / Run / Assert receive d
  no os.Chdir boilerplate

# AssembleRefRootSource
root package: no package free var DOCTEST_ROOT / DOCTEST_SESSION_ID
```

## Preconditions

- Author omits d; leaf path default `nested/leaf`.

## Steps

1. Assemble ref root + leaf.
2. Assert inject contract on leaf; assert no free vars on root.

## Context

- Ref currently mirrors classic Chdir + free vars → RED until P2.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.AuthorDMode = "named-d"
	return nil
}
```
