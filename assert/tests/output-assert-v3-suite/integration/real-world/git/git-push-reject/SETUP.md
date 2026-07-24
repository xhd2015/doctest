# Scenario

**Feature**: git push reject

```
# git push
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"! \\[rejected\\]        main -> main \\(fetch first\\)",
	)
	req.Actual = "! [rejected]        main -> main (fetch first)"
	return nil
}
```
