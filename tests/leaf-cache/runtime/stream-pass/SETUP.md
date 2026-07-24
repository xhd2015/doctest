# Scenario

**Feature**: stream PutPass on each suite-leaf pass so mid-run interrupt keeps already-passed leaves warm

```
# multi-leaf: fast pass + hang leaf (blocks end-of-run record)
doctest test fixture
  -> Action:pass leaf_a -> stream PutPass(key) immediately
  -> hang leaf still running
  -> SIGINT (no clean end-of-run RecordPasses)
# unhang hang leaf (fixture rewrite) then re-run
doctest test fixture -> leaf_a GetPass hit -> Cached >= 1
```

## Preconditions

- Nested CLI integration under `runtime/**` (same `runtime_multi` + isolate env).
- Fixture has one instant-pass leaf and one long-sleep hang leaf so the suite
  cannot reach end-of-run PutPass before interrupt.
- `InterruptAfterDots` kills the nested process after N quiet progress dots.
- After interrupt, mutation rewrites the hang leaf to a fast pass so run2 can finish.

## Steps

1. Parent runtime sets Bin, timeout, isolateRuntimeEnv.
2. Child prepares stream-interrupt fixture and interrupt/mutation knobs.
3. Assert run2 programmatic Cached > 0 for leaves that passed before SIGINT.

## Context

- Without stream PutPass, only end-of-run `recordLeafCachePasses` stores passes;
  SIGINT before hang completes leaves the store cold → run2 Cached == 0 (**RED**
  until product streams PutPass on each countable pass event).
- Fail still never PutPass (covered by `runtime/fail-path/fail-not-stored`).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
