# Scenario

**Feature**: apt-get update

```
# apt-get update
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Reading package lists\\.\\.\\. Done",
	)
	req.Actual = "Reading package lists... Done"
	return nil
}
```
