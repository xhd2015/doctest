# Scenario

**Feature**: cmake

```
# cmake
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"-- Configuring done",
	)
	req.Actual = "-- Configuring done"
	return nil
}
```
