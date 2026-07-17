# Scenario

**Feature**: session-mod cache materialization lifecycle

```
# cold/warm cache under UserCacheDir/doctest/session-mod/<md5>
doctest test with session import -> materialize or reuse
doctest test without session import -> skip materialize
```

## Preconditions

- Cache tests serialize via flock so parallel leaves do not race on wipe/create.
- Leaves may delete the expected cache dir to force cold materialize.

## Steps

1. Acquire cache lock.
2. Leaf prepares module and runs doctest subprocess.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	lockCacheTests(t)
	return nil
}
```
