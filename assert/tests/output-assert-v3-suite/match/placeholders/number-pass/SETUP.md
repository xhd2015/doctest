# Scenario

**Feature**: V3S-M1 — number placeholder pass

```
# Server listen on: __PORT__ matches numeric actual
Matcher <- actual with PORT=8901
```

## Steps
1. Set PORT placeholder template and matching actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("__PORT__: type=number\n", "Server listen on: __PORT__")
	req.Actual = "Server listen on: 8901"
	return nil
}
```