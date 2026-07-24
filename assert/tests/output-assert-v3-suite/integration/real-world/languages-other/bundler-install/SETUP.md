# Scenario

**Feature**: bundle install

```
# bundle install
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Bundle complete! 10 Gemfile dependencies, 20 gems now installed\\.",
	)
	req.Actual = "Bundle complete! 10 Gemfile dependencies, 20 gems now installed."
	return nil
}
```
