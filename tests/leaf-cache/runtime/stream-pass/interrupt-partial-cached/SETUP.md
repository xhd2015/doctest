# Scenario

**Feature**: SIGINT after some pass dots still caches already-passed leaves for the next run

```
run1: 1 fast pass + hang; kill after 1 progress dot (no end-of-run PutPass)
mutate: rewrite hang leaf to fast pass (unhang)
run2: default test -> Cached >= 1 (streamed PutPass survivor)
```

## Preconditions

- `prepareStreamInterruptFixture`: `leaf_a` instant pass; `leaf_hang` sleeps.
- `InterruptAfterDots=1` fires SIGINT after the first quiet progress dot.
- `MutateAfterRun=1` + `Mutation=stream_unhang` so run2 is not blocked by hang.
- Isolated `DOCTEST_LEAF_CACHE`; fresh GOCACHE per invocation.

## Steps

1. Build stream-interrupt fixture.
2. Args = `test <fixture>` with InterruptAfterDots=1.
3. Args2 = `test <fixture>` after unhang mutation.
4. Assert run2 summary Cached >= 1.

## Context

- Proves stream PutPass on JSON `Action: pass` for countable suite leaves.
- End-of-run record may remain as idempotent reconcile; this leaf forbids
  relying on end-of-run alone.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	req.FixtureDir = prepareStreamInterruptFixture(t)
	req.InterruptAfterDots = 1
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", req.FixtureDir}
	req.Mutation = "stream_unhang"
	req.MutateAfterRun = 1
	return nil
}
```
