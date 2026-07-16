# Scenario

**Feature**: psql SELECT

```
# psql
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"",
		" id \\| name \n----\\+------",
	)
	req.Actual = " id | name \n----+------"
	return nil
}
```
