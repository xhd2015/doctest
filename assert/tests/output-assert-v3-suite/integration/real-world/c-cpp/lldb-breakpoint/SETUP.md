# Scenario

**Feature**: lldb

```
# lldb
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__FILE__: 'type=string, example=main.cpp:10'\n",
		"Breakpoint 1: where = main, address = 0x1000, file = 'main\\.cpp', line = 10",
	)
	req.Actual = "Breakpoint 1: where = main, address = 0x1000, file = 'main.cpp', line = 10"
	return nil
}
```
