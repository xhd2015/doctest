# Scenario

**Feature**: wc -l

```
# wc -l
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=42 f.txt'\n",
		"42 f.txt",
	)
	req.Actual = "42 f.txt"
	return nil
}
```
