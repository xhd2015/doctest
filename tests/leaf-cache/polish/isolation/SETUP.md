# Scenario

**Feature**: tree identity mixed into keys (L2 library)

```
# identical relative spine under two abs TreeRoots → distinct ComputeLeafKey digests
```

## Preconditions

- In-process only; unlabeled.
- Twin content-identical trees via `prepareTwinTrees`.

## Steps

1. Child prepares twins and sets Op=`compute_two_inputs`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	return nil
}
```
