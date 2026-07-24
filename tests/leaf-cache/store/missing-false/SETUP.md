# Scenario

**Feature**: GetPass returns false when the key was never written

```
NewStore(emptyRoot) -> st
st.GetPass(neverWritten) -> false, nil
```

## Preconditions

- Fresh StoreRoot with no prior PutPass.
- Key is a synthetic string that was not stored.

## Steps

1. Set Op=`store_missing` and a Key that was never Put.
2. Assert Hit is false and err is nil.

## Context

- Missing is not an error; callers treat false as cache miss.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "store_missing"
	req.Key = "ffffffffffffffffffffffffffffffff"
	return nil
}
```
