# Scenario

**Feature**: pytest

```
# pytest
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__N__: 'type=number, example=3'\n",
		"== 3 passed in 0\\.12s ==",
	)
	req.Actual = "== 3 passed in 0.12s =="
	return nil
}
```
