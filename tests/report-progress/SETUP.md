# Scenario

**Feature**: tests in this group run the `report-progress` binary directly

```
# sub-agents report progress to a file
sub-agent --writes--> progress file (env var DOCTEST_PROGRESS_FILE)

# multiple entries append
each step -> structured JSON entry -> append to file
```

## Preconditions
- Tests in this group run the `report-progress` binary directly.
- The binary is dispatched via the doctest binary copied as `report-progress`.

## Steps
1. Copy the built doctest binary to a temp dir as `report-progress`.
2. Run the binary with the leaf's args and env.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 60 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, "..", ".."))
	req.Env = append(req.Env, "TEST_GROUP=report-progress")
	return nil
}
```
