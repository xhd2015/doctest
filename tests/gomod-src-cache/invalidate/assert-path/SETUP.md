# Scenario

**Feature**: effective assert cache dir change invalidates fingerprint for non-doctest modules

```
first write with WithAssertReplace + AssertCacheDir A
seed tidy-done
second write with AssertCacheDir B (same flags otherwise)
  -> miss rebuild; go.mod gains replace to B; tidy-done dropped if wrote
```

## Steps

1. Use non-doctest ModPath `example.com/app`.
2. First write with assert cache dir under temp path A.
3. Second flags: assert replace still on, cache dir B.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = "example.com/app"
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")

	assertA := filepath.Join(t.TempDir(), "assert-cache-a")
	assertB := filepath.Join(t.TempDir(), "assert-cache-b")
	req.WithAssertReplace = true
	req.AssertCacheDir = assertA

	firstWrite(t, req)
	seedTidyDone(t, req.GenDir)
	req.SnapGoModContentBefore = readFileOrEmpty(filepath.Join(req.GenDir, "go.mod"))
	snapshotGoModMtime(t, req)

	req.UseSecondFlags = true
	req.SecondWithAssertReplace = true
	req.SecondAssertCacheDir = assertB
	req.SecondWithSessionReplace = false
	req.SecondSessionCacheDir = ""
	return nil
}
```
