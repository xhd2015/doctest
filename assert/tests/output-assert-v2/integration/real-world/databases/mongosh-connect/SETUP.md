# Scenario

**Feature**: mongosh

```
# mongosh
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"connecting to: mongodb://127.0.0.1:27017",
	)
	req.Actual = "connecting to: mongodb://127.0.0.1:27017"
	return nil
}
```
