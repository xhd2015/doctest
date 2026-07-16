# Scenario

**Feature**: V3S-M11 — escaped dollar and dots are literals under v3 raw RE

```
# cost: \$5\.00 — $ is not end-anchor; dots are literal
Matcher <- exact cost: $5.00
```

## Steps
1. Set escaped pattern line `cost: \$5\.00` and identical actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "cost: \\$5\\.00")
	req.Actual = "cost: $5.00"
	return nil
}
```
