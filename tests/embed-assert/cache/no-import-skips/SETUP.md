# Scenario

**Feature**: even without author assert import, generation still materializes assert-mod

```
# fixture tree has no github.com/xhd2015/doctest/assert import in SETUP/ASSERT/Run
# but generate always sets assertImport=true → assert-mod cache is created
doctest test -> MaterializeAssertModule for always-on assert-mod
```

## Preconditions

- Module tree does **not** import `github.com/xhd2015/doctest/assert` in author harness.
- Expected cache dir is removed before run so creation is observable (isolated
  `DOCTEST_CACHE_HOME` from parent `cache/SETUP.md`).

## Steps

1. Remove expected assert-mod cache dir.
2. Create public module **without** assert import in fixture sources.
3. Run doctest test (product always materializes assert-mod for external modules).
4. Assert assert-mod cache exists for the content key.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	cacheDir := expectedAssertCacheDir(t)
	_ = os.RemoveAll(cacheDir)
	createPublicModuleProject(t, "", defaultPublicAssertGo())
	setupModuleEnv(t, req)
	req.Args = []string{"test", testDir, "-v"}
	return nil
}
```
