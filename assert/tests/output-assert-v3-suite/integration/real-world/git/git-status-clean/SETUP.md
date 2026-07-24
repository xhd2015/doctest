# Scenario

**Feature**: git status clean

```
# git status
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"",
		"nothing to commit, working tree clean",
	)
	req.Actual = "nothing to commit, working tree clean"
	return nil
}
```
