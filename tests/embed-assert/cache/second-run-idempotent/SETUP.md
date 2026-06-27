# Scenario

**Feature**: second doctest test run does not rewrite existing assert-mod cache files

```
# warm cache after first materialization
doctest test (run 1) -> doctest test (run 2) -> cache bytes and mtimes unchanged
```

## Preconditions

- Assert cache dir is materialized before the second run.

## Steps

1. Create public module with assert leaf.
2. Run `doctest test` once to warm cache.
3. Snapshot cache file mtimes and MD5 digests.
4. Run `doctest test` again with same tree.

```go
import (
	"path/filepath"
	"testing"
)

var (
	beforeAssertMtime	int64
	beforeAssertDigest	[16]byte
	beforeGoModMtime	int64
	beforeGoModDigest	[16]byte
)

func Setup(t *testing.T, req *Request) error {
	createPublicModuleProject(t, "", defaultAssertAssertGo())
	setupModuleEnv(t, req)

	req.Args = []string{"test", testDir, "-v"}
	if _, err := Run(t, req); err != nil {
		t.Fatalf("first doctest run failed: %v", err)
	}

	cacheDir := expectedAssertCacheDir(t)
	assertCacheLayout(t, cacheDir)
	beforeAssertMtime, beforeAssertDigest = snapshotFileState(t, filepath.Join(cacheDir, "assert.go"))
	beforeGoModMtime, beforeGoModDigest = snapshotFileState(t, filepath.Join(cacheDir, "go.mod"))

	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```