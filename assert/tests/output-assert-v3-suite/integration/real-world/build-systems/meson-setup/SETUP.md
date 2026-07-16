# Scenario

**Feature**: meson setup

```
# meson setup
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__DIR__: 'type=string, example=build'\n",
		"The Meson build system\nDirectory build created",
	)
	req.Actual = "The Meson build system\nDirectory build created"
	return nil
}
```
