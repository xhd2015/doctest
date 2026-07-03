# Scenario

**Feature**: gcloud auth

```
# gcloud auth list
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"",
		"Credentialed Accounts",
	)
	req.Actual = "Credentialed Accounts"
	return nil
}
```
