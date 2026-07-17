# Scenario

**Feature**: two exclusive creates in the same second produce distinct paths

```
# clash avoidance
CreateRunFile(same second) x2 -> path1 != path2; both exist as files
```

## Preconditions

- Same `At` timestamp and same optional suffix seed for both creates.
- Implementation may vary `NN` (00–99) and/or generate a unique suffix.

## Steps

1. Call `CreateRunFile` twice with identical clock input.
2. Compare returned paths and ensure both files exist (empty ok).

## Context

- Exclusive create (O_EXCL or equivalent) prevents clobbering.
- Either different `NN` or different suffix is acceptable; paths must differ.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "create_run_files"
	req.CacheDir = t.TempDir()
	req.ProjectID = "github.com_xhd2015_doctest"
	req.At = time.Date(2026, 7, 17, 15, 0, 0, 0, time.UTC)
	req.Suffix = "" // implementation may generate unique suffixes
	req.CreateCount = 2
	return nil
}
```
