# Scenario

**Feature**: git push reject

```
# git push
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__LINE__: 'type=string, example=! [rejected]        main -> main (fetch first)'\n",
		"__LINE__",
	)
	req.Actual = "! [rejected]        main -> main (fetch first)"
	return nil
}
```
