# Scenario

**Feature**: v2 `<ansi-color>` ANSI span matching

```
# <ansi-color SPEC>inner</ansi-color> asserts strict SGR envelope
Matcher compares literal text wrapped in expected ANSI codes
```

## Steps
1. Templates use `<ansi-color>` (same tag name as v1).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```