# Scenario

**Feature**: composer install

```
# composer install
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Generating autoload files",
	)
	req.Actual = "Generating autoload files"
	return nil
}
```
