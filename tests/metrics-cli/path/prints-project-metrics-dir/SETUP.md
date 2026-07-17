# Scenario

**Feature**: path prints metrics dir when the directory already exists

```
mkdir $root/doctest/metrics/<id>
doctest metrics path -> prints that directory
```

## Preconditions

- Metrics project directory exists (empty runs ok).

## Steps

1. Create project metrics directory under MetricsRoot.
2. Run `metrics path`.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	ensureFixtureProject(t, req)
	if err := os.MkdirAll(projectMetricsDir(req), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	req.Args = []string{"metrics", "path"}
	return nil
}
```
