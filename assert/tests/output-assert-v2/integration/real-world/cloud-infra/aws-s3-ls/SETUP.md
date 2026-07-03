# Scenario

**Feature**: aws s3 ls

```
# aws s3 ls
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__BUCKET__: 'type=string, example=my-bucket'\n",
		"PRE my-bucket/",
	)
	req.Actual = "PRE my-bucket/"
	return nil
}
```
