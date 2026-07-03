# Scenario

**Feature**: cargo bench

```
# cargo bench
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__NS__: 'type=number, example=1000'\n",
		"bench: 1000 ns/iter",
	)
	req.Actual = "bench: 1000 ns/iter"
	return nil
}
```
