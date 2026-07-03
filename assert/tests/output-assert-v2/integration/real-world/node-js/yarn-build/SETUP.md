# Scenario

**Feature**: yarn build

```
# yarn build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Done in 1.23s.",
	)
	req.Actual = "Done in 1.23s."
	return nil
}
```
