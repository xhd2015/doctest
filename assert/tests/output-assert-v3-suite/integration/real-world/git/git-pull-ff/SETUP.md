# Scenario

**Feature**: git pull

```
# git pull
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Fast-forward",
	)
	req.Actual = "Fast-forward"
	return nil
}
```
