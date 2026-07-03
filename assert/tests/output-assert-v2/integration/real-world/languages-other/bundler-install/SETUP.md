# Scenario

**Feature**: bundle install

```
# bundle install
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Bundle complete! 10 Gemfile dependencies, 20 gems now installed.",
	)
	req.Actual = "Bundle complete! 10 Gemfile dependencies, 20 gems now installed."
	return nil
}
```
