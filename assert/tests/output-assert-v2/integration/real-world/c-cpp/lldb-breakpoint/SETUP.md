# Scenario

**Feature**: lldb

```
# lldb
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__FILE__: 'type=string, example=main.cpp:10'\n",
		"Breakpoint 1: where = main, address = 0x1000, file = 'main.cpp', line = 10",
	)
	req.Actual = "Breakpoint 1: where = main, address = 0x1000, file = 'main.cpp', line = 10"
	return nil
}
```
