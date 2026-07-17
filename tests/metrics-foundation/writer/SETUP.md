# Scenario

**Feature**: buffered JSONL writer flushes at 128KiB and on Close

```
# small write
Write(small) -> buffer only -> file empty or unchanged until Close -> flush all

# large write
Write(≥128KiB cumulative) -> mid-run flush -> file non-empty before Close

# partial
Write(run_start, leaf_*) without run_end -> Close -> lines still readable
```

## Preconditions

- Temp cache; `FlushThreshold` is 128 * 1024 bytes.
- One process owns one run file.

## Steps

1. Leaf configures small, large, or partial sequences.
2. Run opens writer, writes, optionally stats before close, closes, reads.
3. Assert sizes and decodable lines.

## Context

- No fsync required.
- Append-only JSON lines with trailing `\n`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "write_sequence"
	req.CacheDir = t.TempDir()
	req.ProjectID = "writer_proj"
	req.Suffix = "wr000001"
	return nil
}
```
