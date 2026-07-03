# Scenario

**Feature**: go test -cover

```
# go test -cover
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__PCT__: 'type=number, example=80.5'\n",
		"coverage: 80.5% of statements",
	)
	req.Actual = "coverage: 80.5% of statements"
	return nil
}
```
