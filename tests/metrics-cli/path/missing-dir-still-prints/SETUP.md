# Scenario

**Feature**: path still prints canonical location when metrics dir is missing

```
# MetricsRoot empty — no doctest/metrics/<id> yet
doctest metrics path -> exit 0, print canonical path
```

## Preconditions

- Do not create the project metrics directory.

## Steps

1. Prepare WorkDir + empty MetricsRoot only.
2. Run `metrics path`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	// intentionally do not MkdirAll project metrics dir
	req.Args = []string{"metrics", "path"}
	return nil
}
```
