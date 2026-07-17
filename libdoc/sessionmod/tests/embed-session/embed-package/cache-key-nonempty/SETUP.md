# Scenario

**Feature**: RawSourceCacheKeyMD5 is a stable non-empty cache directory key

```
RawSourceCacheKeyMD5() -> non-empty hex used as session-mod/<md5>/
```

## Steps

1. Call RawSourceCacheKeyMD5 and Content.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	runKind = "sessionmod-cache-key"
	return nil
}
```
