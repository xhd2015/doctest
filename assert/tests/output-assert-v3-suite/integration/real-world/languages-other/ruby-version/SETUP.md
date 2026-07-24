# Scenario

**Feature**: ruby -v

```
# ruby
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__VER__: 'type=string, example=3.2.0'\n",
		"ruby 3\\.2\\.0",
	)
	req.Actual = "ruby 3.2.0"
	return nil
}
```
