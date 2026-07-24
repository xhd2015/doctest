# Scenario

**Feature**: helm install

```
# helm install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__REL__: 'type=string, example=myrel'\n",
		"NAME: myrel",
	)
	req.Actual = "NAME: myrel"
	return nil
}
```
