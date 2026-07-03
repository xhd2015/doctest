# Scenario

**Feature**: V2-M11 — mid-line dollar stays pattern line

```
# cost: $5.00 has $ not at line end — literal, not end anchor
Matcher <- exact cost: $5.00
```

## Steps
1. Set pattern line `cost: $5.00` and identical actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "cost: $5.00")
	req.Actual = "cost: $5.00"
	return nil
}
```