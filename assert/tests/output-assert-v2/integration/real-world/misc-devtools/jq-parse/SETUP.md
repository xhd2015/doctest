# Scenario

**Feature**: jq

```
# jq .name
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"\"alice\"",
	)
	req.Actual = "\"alice\""
	return nil
}
```
