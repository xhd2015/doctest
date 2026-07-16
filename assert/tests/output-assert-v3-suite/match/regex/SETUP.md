# Scenario

**Feature**: Intentional raw-RE template lines under v3

```
# content lines are always raw Go regexp; these leaves use RE deliberately
Matcher <- actual line must fully match regex
```

## Steps
1. Body lines use deliberate RE (e.g. `.*`, alternation) — not escaped as literals.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	return nil
}
```
