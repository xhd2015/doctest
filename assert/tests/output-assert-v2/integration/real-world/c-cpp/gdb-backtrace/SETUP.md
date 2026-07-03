# Scenario

**Feature**: gdb bt

```
# gdb
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"#0  main () at main.c:5",
	)
	req.Actual = "#0  main () at main.c:5"
	return nil
}
```
