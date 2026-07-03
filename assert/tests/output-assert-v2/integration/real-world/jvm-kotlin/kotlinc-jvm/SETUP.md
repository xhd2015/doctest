# Scenario

**Feature**: kotlinc -jvm

```
# kotlinc
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__VER__: 'type=string, example=21'\n",
		"jvm target 21",
	)
	req.Actual = "jvm target 21"
	return nil
}
```
