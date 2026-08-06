# Scenario

**Feature**: gen-plan plan phase prints arg hierarchies on stderr before/during generate

```
DOCTEST_DEBUG=gen-plan=1,bypass-go-test=1
  + doctest test <fixture> --gen-dir <isolated>
  -> stderr: gen-plan: invocation / arg[i/n] / [merged]
```

## Preconditions

- Product CLI leaves (`label: e2e` on descendants).
- Always combine `bypass-go-test=1` so go test is skipped (fast plan/result).
- Isolated GenDir + cache home per leaf.

## Steps

1. Group Setup sets Mode=cli and DebugEnv combo.
2. Child leaves prepare single- or multi-arg fixtures and Args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Mode = "cli"
	req.DebugEnv = debugGenPlanBypass
	return nil
}
```
