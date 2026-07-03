# Scenario

**Feature**: npm install

```
# npm install
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__N__: 'type=number, example=42'\n",
		"added 42 packages, and audited 43 packages in 2s",
	)
	req.Actual = "added 42 packages, and audited 43 packages in 2s"
	return nil
}
```
