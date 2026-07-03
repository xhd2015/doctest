# Scenario

**Feature**: V2-M2 — number placeholder fail

```
# non-numeric value rejected for type=number
Matcher <- actual with PORT=abc
```

## Steps
1. Same template as M1 with non-numeric actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("__PORT__: type=number\n", "Server listen on: __PORT__")
	req.Actual = "Server listen on: abc"
	return nil
}
```