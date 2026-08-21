# Scenario

**Feature**: installing a versioned skill preserves its complete document

```
doctest skill dev-test --install <temp-dir> -> <temp-dir>/SKILL.md
```

## Steps

1. Install `dev-test` to an explicit temporary skill directory.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.InstallDir = t.TempDir()
	req.Args = []string{"skill", "dev-test", "--install", req.InstallDir}
	return nil
}
```
