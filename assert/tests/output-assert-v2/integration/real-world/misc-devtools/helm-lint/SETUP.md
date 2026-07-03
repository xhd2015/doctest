# Scenario

**Feature**: helm lint

```
# helm lint
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"==> Linting chart\n...1 lines omitted...\n1 chart(s) linted, 0 chart(s) failed",
	)
	req.Actual = "==> Linting chart\n[INFO] Chart.yaml: icon is recommended\n1 chart(s) linted, 0 chart(s) failed"
	return nil
}
```
