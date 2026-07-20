# Scenario

**Bug**: multi-tree prepare must not rewrite go.mod mtime for doctest self-module

```
WriteGoMod(assert=false)
  -> seed tidy-done, force old go.mod mtime
WriteGoMod(assert=true, session=true, cache paths)
  -> go.mod mtime unchanged
  -> tidy-done retained
  -> still no gomod-fp; manifest present
```

## Steps

1. First WriteGoMod without assert/session replaces.
2. Seed tidy-done; snapshot go.mod mtime.
3. Second WriteGoMod with assert+session flags (ineffective for doctest module).

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Mode = "write-gomod-second"
	req.WithAssertReplace = false
	req.AssertCacheDir = ""
	req.WithSessionReplace = false
	req.SessionCacheDir = ""
	req.SecondWithAssertReplace = true
	req.SecondAssertCacheDir = "/tmp/assert-cache-gen-write-manifest"
	req.SecondWithSessionReplace = true
	req.SecondSessionCacheDir = "/tmp/session-cache-gen-write-manifest"
	prepareFreshGen(t, req, "module github.com/xhd2015/doctest\n\ngo 1.21\n")
	firstWriteGoMod(t, req)
	seedTidyDone(t, req.GenDir)
	snapshotGoModMtime(t, req)
	snapshotManifestMtime(t, req)
	return nil
}
```
