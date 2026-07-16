# Scenario

**Feature**: sqlite3

```
# sqlite3
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__N__: 'type=number, example=1'\n",
		"1",
	)
	req.Actual = "1"
	return nil
}
```
