# Scenario

**Feature**: gcc compile

```
# gcc -c
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"cc -c main.c",
	)
	req.Actual = "cc -c main.c"
	return nil
}
```
