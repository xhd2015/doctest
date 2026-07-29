# Scenario

**Feature**: auto + no xgo mock in entry graph resolves to `go`

```
modRoot: example.com/app/runpkg (no mock imports)
  -> DetectXgoMockUsage -> needsXgo=false
  -> ResolveGoTestCmd(auto, false) -> "go"
```

## Preconditions

- Fixture has a single entry package with no import of
  `github.com/xhd2015/xgo/runtime/mock` (direct or transitive under the module).

## Steps

1. Seed no-mock module; entry import path `example.com/app/runpkg`.
2. Run detect + resolve (auto).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	seedNoMockModule(t, req)
	req.GoCmdFlag = "" // auto / omitted
	return nil
}
```
