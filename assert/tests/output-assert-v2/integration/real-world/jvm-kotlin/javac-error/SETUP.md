# Scenario

**Feature**: javac error

```
# javac
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Main.java:3: error: cannot find symbol",
	)
	req.Actual = "Main.java:3: error: cannot find symbol"
	return nil
}
```
