# Scenario

**Feature**: cmake

```
# cmake
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"-- Configuring done",
	)
	req.Actual = "-- Configuring done"
	return nil
}
```
