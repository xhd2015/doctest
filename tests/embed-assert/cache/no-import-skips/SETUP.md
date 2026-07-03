# Scenario

**Feature**: doctest run without assert import does not create new assert-mod cache entry

```
# no assert import in tree
doctest test -> assert-mod cache entry count unchanged
```

## Preconditions

- Leaf does not import `github.com/xhd2015/doctest/assert`.
- Snapshot assert-mod cache entries before run.

## Steps

1. Record whether the current `RawSourceCacheKeyMD5` cache dir exists.
2. Create public module without assert import.
3. Run `doctest test <tests> -v`.

```go
import (
	"os"
	"testing"
)

var cacheDirExistedBefore bool

func Setup(t *testing.T, req *Request) error {
	cacheDir := expectedAssertCacheDir(t)
	_, err := os.Stat(cacheDir)
	cacheDirExistedBefore = err == nil
	createPublicModuleProject(t, "", defaultPublicAssertGo())
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```