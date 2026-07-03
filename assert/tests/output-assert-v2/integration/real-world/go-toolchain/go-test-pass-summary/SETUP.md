# Scenario

**Feature**: go test pass

```
# go test
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__SEC__: 'type=number, example=0.12'\n",
		"PASS\nok  \texample.com/x\t0.12s",
	)
	req.Actual = "PASS\nok  \texample.com/x\t0.12s"
	return nil
}
```
