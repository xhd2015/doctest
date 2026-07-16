# Scenario

**Feature**: yarn build

```
# yarn build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Done in 1\\.23s\\.",
	)
	req.Actual = "Done in 1.23s."
	return nil
}
```
