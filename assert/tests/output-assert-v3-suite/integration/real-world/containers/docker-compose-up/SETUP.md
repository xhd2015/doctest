# Scenario

**Feature**: docker compose up

```
# compose up
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		" Container app  Started",
	)
	req.Actual = " Container app  Started"
	return nil
}
```
