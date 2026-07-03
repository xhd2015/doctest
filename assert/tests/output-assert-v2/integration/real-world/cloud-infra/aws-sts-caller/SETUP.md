# Scenario

**Feature**: aws sts

```
# aws sts get-caller-identity
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__ARN__: 'type=string, example=arn:aws:iam::1:user/u'\n",
		"arn:aws:iam::1:user/u",
	)
	req.Actual = "arn:aws:iam::1:user/u"
	return nil
}
```
