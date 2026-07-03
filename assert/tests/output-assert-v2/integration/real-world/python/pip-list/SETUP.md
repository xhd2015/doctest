# Scenario

**Feature**: pip list

```
# pip list
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__PKG__: 'type=string, example=requests'\n",
		"requests                2.31.0",
	)
	req.Actual = "requests                2.31.0"
	return nil
}
```
