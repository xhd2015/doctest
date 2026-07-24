# Scenario

**Feature**: RawSourceCacheKeyMD5 is a stable non-empty cache directory key

```
RawSourceCacheKeyMD5() -> non-empty hex used as session-mod/<md5>/
```

## Steps

1. Call RawSourceCacheKeyMD5 and Content.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.RunKind = "sessionmod-cache-key"
	return nil
}
```
