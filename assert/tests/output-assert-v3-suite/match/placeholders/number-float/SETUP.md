# Scenario

**Feature**: V3S-M14 — float number placeholder pass

```
# type=number accepts floats (-?\d+(\.\d+)?)
Matcher <- actual latency: 3.14
```

## Steps
1. Set MS number placeholder and float actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("__MS__: type=number\n", "latency: __MS__")
	req.Actual = "latency: 3.14"
	return nil
}
```