# Scenario

**Feature**: composer install

```
# composer install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		"Generating autoload files",
	)
	req.Actual = "Generating autoload files"
	return nil
}
```
