# Scenario

**Feature**: cargo add

```
# cargo add
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__CRATE__: 'type=string, example=serde'\n",
		"      Adding serde v1\\.0\\.0 to dependencies",
	)
	req.Actual = "      Adding serde v1.0.0 to dependencies"
	return nil
}
```
