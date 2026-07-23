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
- **L3 e2e**: product binary (UserCacheDir side effects / nested test).

## Steps

1. Acquire cache lock; set UseCLI + Bin.
2. Leaf prepares module and runs doctest subprocess.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	lockCacheTests(t)
	req.UseCLI = true
	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	return nil
}
```
