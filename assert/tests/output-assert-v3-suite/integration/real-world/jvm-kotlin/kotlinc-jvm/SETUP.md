# Scenario

**Feature**: kotlinc -jvm

```
# kotlinc
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__VER__: 'type=string, example=21'\n",
		"jvm target 21",
	)
	req.Actual = "jvm target 21"
	return nil
}
```
