# Scenario

**Feature**: PutPass then GetPass returns true for the same key

```
NewStore(root) -> st
st.PutPass(key)
st.GetPass(key) -> true
```

## Preconditions

- Isolated temp StoreRoot from parent.
- Fixed synthetic key (no fixture required).

## Steps

1. Set Op=`store_put_get` and a non-empty Key.
2. Run PutPass then GetPass.
3. Assert Hit is true.

## Context

- Only explicit PutPass creates a pass marker.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "store_put_get"
	req.Key = "aabbccddeeff00112233445566778899"
	return nil
}
```
