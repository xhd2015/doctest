# Scenario

**Feature**: helm install

```
# helm install
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__REL__: 'type=string, example=myrel'\n",
		"NAME: myrel",
	)
	req.Actual = "NAME: myrel"
	return nil
}
```
