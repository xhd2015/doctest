# Scenario

**Feature**: sed

```
# sed s///
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"hello earth",
	)
	req.Actual = "hello earth"
	return nil
}
```
