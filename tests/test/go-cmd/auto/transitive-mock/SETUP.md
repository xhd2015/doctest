# Scenario

**Feature**: auto + **transitive** xgo mock (Run → project helper → mock) → `xgo`

```
runpkg.Run imports example.com/app/helper
helper imports github.com/xhd2015/xgo/runtime/mock
  -> DetectXgoMockUsage([example.com/app/runpkg]) -> needsXgo=true
  -> ResolveGoTestCmd(auto, true) -> "xgo"
```

## Preconditions

- Entry package must **not** import mock directly; only `helper` does.
- Detection must walk project imports (credit-pricing-center style), not only
  markdown AST for the string `mock.Patch`.

## Steps

1. Seed transitive mock module.
2. Run detect + resolve (auto).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	seedTransitiveMockModule(t, req)
	req.GoCmdFlag = "" // auto
	return nil
}
```
