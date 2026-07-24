# Scenario

**Feature**: V3S-M2 — number placeholder fail

```
# non-numeric value rejected for type=number
Matcher <- actual with PORT=abc
```

## Steps
1. Same template as M1 with non-numeric actual.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__PORT__: type=number\n", "Server listen on: __PORT__")
	req.Actual = "Server listen on: abc"
	return nil
}
```