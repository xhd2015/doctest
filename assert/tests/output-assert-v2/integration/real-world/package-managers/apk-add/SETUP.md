# Scenario

**Feature**: apk add

```
# apk add
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__PKG__: 'type=string, example=git'\n",
		"OK: 50 MiB in 20 packages",
	)
	req.Actual = "OK: 50 MiB in 20 packages"
	return nil
}
```
