# Scenario

**Feature**: zig build

```
# zig build
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"install",
	)
	req.Actual = "install"
	return nil
}
```
