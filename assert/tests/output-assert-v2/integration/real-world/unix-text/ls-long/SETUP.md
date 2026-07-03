# Scenario

**Feature**: ls -l

```
# ls -l
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=-rw-r--r-- main.go'\n",
		"-rw-r--r-- main.go",
	)
	req.Actual = "-rw-r--r-- main.go"
	return nil
}
```
