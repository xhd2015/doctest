# Scenario

**Feature**: compile strategy determines how generated tests resolve the assert module

```
# legacy nested module: no internal/ imports
doctest -> WriteGoMod testcase -> replace parent + replace assert => cache

# internal compile: parent internal/ imports detected
doctest -> .doctest_run_* -> -modfile (parent go.mod + assert replace)
```

## Preconditions

- Siblings split on compile strategy: nested-module vs internal-compile.
- Both paths require `CasesImportAssertPackage` for assert-specific wiring.

## Steps

1. Descendant Setup selects fixture and doctest args for its strategy.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// No process Env (Parallel-safe). Temp fixtures have no go.work.
	_ = req
	return nil
}
```