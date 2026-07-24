# Scenario

**Feature**: `metrics path` prints the project metrics directory

```
cwd origin + DOCTEST_METRICS_ROOT
  -> doctest metrics path
  -> stdout: $root/doctest/metrics/<project_id>
```

## Preconditions

- Project identity comes from FixtureOrigin in WorkDir.
- Path is absolute under MetricsRoot.

## Steps

1. Seed git WorkDir with fixture origin.
2. Run `metrics path`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureFixtureProject(t, req)
	return nil
}
```
