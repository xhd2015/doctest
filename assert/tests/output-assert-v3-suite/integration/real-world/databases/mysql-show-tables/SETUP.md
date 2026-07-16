# Scenario

**Feature**: mysql SHOW TABLES

```
# mysql
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__TBL__: 'type=string, example=users'\n",
		"users",
	)
	req.Actual = "users"
	return nil
}
```
