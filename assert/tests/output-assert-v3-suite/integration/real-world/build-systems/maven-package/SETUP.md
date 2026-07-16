# Scenario

**Feature**: mvn package

```
# mvn package
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"\\[INFO\\] BUILD SUCCESS",
	)
	req.Actual = "[INFO] BUILD SUCCESS"
	return nil
}
```
