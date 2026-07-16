# Scenario

**Feature**: npm ci

```
# npm ci
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"added 100 packages in 5s",
	)
	req.Actual = "added 100 packages in 5s"
	return nil
}
```
