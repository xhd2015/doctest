# Scenario

**Feature**: mongosh

```
# mongosh
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"connecting to: mongodb://127\\.0\\.0\\.1:27017",
	)
	req.Actual = "connecting to: mongodb://127.0.0.1:27017"
	return nil
}
```
