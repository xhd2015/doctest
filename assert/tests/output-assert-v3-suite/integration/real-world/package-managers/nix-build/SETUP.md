# Scenario

**Feature**: nix-build

```
# nix-build
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PATH__: 'type=string, example=/nix/store/abc'\n",
		"/nix/store/abc",
	)
	req.Actual = "/nix/store/abc"
	return nil
}
```
