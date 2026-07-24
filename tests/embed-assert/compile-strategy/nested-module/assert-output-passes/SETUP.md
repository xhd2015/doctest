# Scenario

**Feature**: generated test using assert.Output compiles and passes via nested module replace

```
# assert.Output in ASSERT.md
doctest test -> replace assert => cache -> go test PASS
```

## Preconditions

- Leaf ASSERT calls `assert.Output` with matching template and actual strings.

## Steps

1. Create public module with assert.Output leaf.
2. Run `doctest test <tests> -v` (mapping-gen cache path).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createPublicModuleProject(t, req, "", defaultAssertAssertGo())
	setupModuleEnv(t, req)
	req.Args = []string{"test", req.TestDir, "-v"}
	return nil
}
```