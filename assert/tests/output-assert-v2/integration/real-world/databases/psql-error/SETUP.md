# Scenario

**Feature**: psql error

```
# psql
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"ERROR:  relation \"x\" does not exist",
	)
	req.Actual = "ERROR:  relation \"x\" does not exist"
	return nil
}
```
