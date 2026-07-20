# Scenario

**Feature**: cacheDir passed to fn is under t.TempDir()/session-once/... and is writable

```
# Once allocates cacheDir under test temp (not UserCacheDir)
session.Once -> cacheDir = t.TempDir()/session-once/<slug>/
fn writes probe file into cacheDir
Caller <- JSON containing path; disk has probe-write
```

## Steps

1. Mode `cache-probe` writes a marker file inside cacheDir.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = "cache-probe"
	req.Key = "layout-key"
	req.SessionID = "once-doctest-layout-" + t.Name()
	req.CallTwice = false
	return nil
}
```
