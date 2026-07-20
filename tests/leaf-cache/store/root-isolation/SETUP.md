# Scenario

**Feature**: pass markers are isolated per store Root

```
stA := NewStore(rootA)
stB := NewStore(rootB)
stA.PutPass(key)
stA.GetPass(key) -> true
stB.GetPass(key) -> false
```

## Preconditions

- Two distinct temp roots: StoreRoot and StoreRootB.
- Same key string used on both stores.

## Steps

1. Set Op=`store_isolate`.
2. PutPass on A only.
3. Assert Hit true on A and HitB false on B.

## Context

- Tests never share the user `$CacheHome`; isolation proves Root is authoritative.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "store_isolate"
	req.Key = "11223344556677889900aabbccddeeff"
	base := t.TempDir()
	req.StoreRoot = filepath.Join(base, "a")
	req.StoreRootB = filepath.Join(base, "b")
	return nil
}
```
