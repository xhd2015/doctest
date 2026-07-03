# Scenario

**Feature**: mvn package

```
# mvn package
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=[INFO] BUILD SUCCESS'\n",
		"__LINE__",
	)
	req.Actual = "[INFO] BUILD SUCCESS"
	return nil
}
```
