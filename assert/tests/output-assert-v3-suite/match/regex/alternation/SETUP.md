# Scenario

**Feature**: V3S-M13 — regex alternation pass

```
# (ok|fail) alternation signal classifies line as regex
Matcher <- actual ok
```

## Steps
1. Set alternation regex line and matching branch `ok`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "(ok|fail)")
	req.Actual = "ok"
	return nil
}
```