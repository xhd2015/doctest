# Scenario

**Feature**: mvn test

```
# mvn test
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Tests run: 3, Failures: 0, Errors: 0, Skipped: 0",
	)
	req.Actual = "Tests run: 3, Failures: 0, Errors: 0, Skipped: 0"
	return nil
}
```
