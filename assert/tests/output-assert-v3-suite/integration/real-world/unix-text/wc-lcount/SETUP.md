# Scenario

**Feature**: wc -l

```
# wc -l
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__LINE__: 'type=string, example=42 f.txt'\n",
		"42 f\\.txt",
	)
	req.Actual = "42 f.txt"
	return nil
}
```
