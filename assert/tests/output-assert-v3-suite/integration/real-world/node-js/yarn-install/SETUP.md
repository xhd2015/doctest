# Scenario

**Feature**: yarn install

```
# yarn
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"success Saved lockfile\\.",
	)
	req.Actual = "success Saved lockfile."
	return nil
}
```
