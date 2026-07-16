# Scenario

**Feature**: php -v

```
# php
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"PHP 8\\.2\\.0 \\(cli\\)",
	)
	req.Actual = "PHP 8.2.0 (cli)"
	return nil
}
```
