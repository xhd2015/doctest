# Scenario

**Feature**: auto cold gen root is `$CacheHome/doctest/mapping-gen-cold`

```
# auto (no --gen-dir)
doctest test --cold-cache <tiny-tree>
  + DOCTEST_CACHE_HOME=<sandbox>
  -> coldHome = $CacheHome/doctest/mapping-gen-cold
  -> RemoveAll+MkdirAll(coldHome) on startup
  -> generate into coldHome; leave content after finish
  -> stderr announces cold-cache mode
```

## Preconditions

- `--gen-dir` is omitted.
- Sandbox `DOCTEST_CACHE_HOME` is set; cold home may be pre-seeded with a marker file.

## Steps

1. Create cache sandbox and tiny fixture project.
2. Seed `marker-before` under cold home (proves startup wipe).
3. Run `doctest test --cold-cache <testDir>` (no `--gen-dir`, no `-count`).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	withCacheSandbox(t, req)
	req.CCTestDir = createTempTestProject(t)
	seedMarker(t, req, req.CCColdHome, "marker-before")
	req.Args = []string{"test", "--cold-cache", req.CCTestDir}
	return nil
}
```
