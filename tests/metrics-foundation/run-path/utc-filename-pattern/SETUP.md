# Scenario

**Feature**: RunFilePath joins cache layout and UTC filename components

```
# fixed instant UTC
2026-07-17 12:34:56 UTC, nn=7, suffix=abc12def
  -> .../runs/2026-07-17-12-34-56-07-abc12def.jsonl
```

## Preconditions

- Cache dir, project id, fixed UTC time, NN, and suffix are provided.

## Steps

1. Call `RunFilePath` with known inputs.
2. Inspect absolute path segments and basename.

## Context

- `NN` is zero-padded to two digits (`07`).
- Directory segments: `<cache>/doctest/metrics/<project_id>/runs/`.

```go
import (
	"path/filepath"
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_file_path"
	req.CacheDir = filepath.Join(t.TempDir(), "cache")
	req.ProjectID = "github.com_xhd2015_doctest"
	req.At = time.Date(2026, 7, 17, 12, 34, 56, 0, time.UTC)
	req.NN = 7
	req.Suffix = "abc12def"
	return nil
}
```
