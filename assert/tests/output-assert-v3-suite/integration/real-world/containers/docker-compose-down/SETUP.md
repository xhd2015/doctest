# Scenario

**Feature**: docker compose down

```
# compose down
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		" Network default  Removed",
	)
	req.Actual = " Network default  Removed"
	return nil
}
```
