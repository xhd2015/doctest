# Scenario

**Feature**: cacheDir passed to fn is under sessions/.../once-... and is writable

```
# Once allocates cacheDir
session.Once -> cacheDir = $UserCacheDir/doctest/sessions/<sid>/once-<slug>/
fn writes probe file into cacheDir
Caller <- JSON containing path; disk has probe-write
```

## Preconditions

- Mode `cache-probe` writes a marker file inside cacheDir.
- Session id and key produce a stable slug path.

## Steps

1. Call Once once with cache-probe mode.
2. Assert CacheDir path prefix/shape and marker file exists.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.SessionID = "once-doctest-layout-" + DOCTEST_SESSION_ID
	req.Key = "layout-probe"
	req.Mode = "cache-probe"
	req.CallTwice = false
	return nil
}
```
