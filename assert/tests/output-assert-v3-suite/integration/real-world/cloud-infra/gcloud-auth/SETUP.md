# Scenario

**Feature**: gcloud auth

```
# gcloud auth list
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"Credentialed Accounts",
	)
	req.Actual = "Credentialed Accounts"
	return nil
}
```
