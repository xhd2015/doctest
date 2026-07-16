# Scenario

**Feature**: sentry-cli

```
# sentry-cli upload
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"> Uploaded release files to Sentry",
	)
	req.Actual = "> Uploaded release files to Sentry"
	return nil
}
```
