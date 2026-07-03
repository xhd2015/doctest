# Scenario

**Feature**: php -v

```
# php
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"PHP 8.2.0 (cli)",
	)
	req.Actual = "PHP 8.2.0 (cli)"
	return nil
}
```
