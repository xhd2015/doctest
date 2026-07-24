# Scenario

**Feature**: metrics foundation APIs for project identity, run JSONL paths, and buffered event writing

```
# resolve identity
git origin URL | abs root -> project_id slug

# allocate run file under cache
$CACHE/doctest/metrics/<project_id>/runs/YYYY-MM-DD-HH-MM-SS-NN-<suffix>.jsonl

# append schema_version 1 events through 128KiB buffered writer
Caller -> Writer.Write(event) -> buffer -> flush (≥128KiB | Close) -> JSONL file
```

## Preconditions

- Package under test: `github.com/xhd2015/doctest/libdoc/metrics` (created by implementer).
- Leaves set `req.Op` and operation-specific fields; root does not invent scenarios.
- Temp cache roots use `t.TempDir()` so tests never write into the user cache.

## Steps

1. Leaf `Setup` fills `Request` for one MECE branch (identity, path, events, or writer).
2. Root `Run` dispatches on `req.Op` to the corresponding metrics API.
3. Leaf `Assert` checks slugs, path patterns, decoded JSONL, or flush sizes.

## Context

- Expected public API surface (implementer contract):
  - `const SchemaVersion = 1`
  - `const FlushThreshold = 128 * 1024`
  - `ProjectIDFromOrigin(origin string) string`
  - `ProjectIDFallback(absRoot string) string` → `nogit_` + first 12 hex chars of SHA-256(absRoot)
  - `RunFilePath(cacheDir, projectID string, t time.Time, nn int, suffix string) string`
  - `CreateRunFile(cacheDir, projectID string, at time.Time, suffix string) (path string, error)` exclusive create with NN disambiguation
  - `OpenWriter(path string) (*Writer, error)`
  - `(*Writer).Write(v any) error` — one JSON object + `\n`, buffered
  - `(*Writer).Close() error` — flush remaining buffer
- No dirty-git field on events. Skip leaves may be recorded as `leaf_end` only.
- Out of scope: CLI wiring, WARNING, skill, `doctest metrics`.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Defaults shared by all leaves; leaves override Op and fields.
	if req.At.IsZero() {
		// Fixed UTC instant for deterministic path tests when leaves do not set At.
		req.At = time.Date(2026, 7, 17, 12, 34, 56, 0, time.UTC)
	}
	if !req.LeaveOpen {
		req.CloseWriter = true
	}
	return nil
}
```
