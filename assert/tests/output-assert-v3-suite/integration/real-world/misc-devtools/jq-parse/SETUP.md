# Scenario

**Feature**: jq

```
# jq .name
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"\"alice\"",
	)
	req.Actual = "\"alice\""
	return nil
}
```
