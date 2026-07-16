# Scenario

**Feature**: gcc compile

```
# gcc -c
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"cc -c main\\.c",
	)
	req.Actual = "cc -c main.c"
	return nil
}
```
