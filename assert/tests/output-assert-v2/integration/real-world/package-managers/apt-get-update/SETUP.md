# Scenario

**Feature**: apt-get update

```
# apt-get update
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Reading package lists... Done",
	)
	req.Actual = "Reading package lists... Done"
	return nil
}
```
