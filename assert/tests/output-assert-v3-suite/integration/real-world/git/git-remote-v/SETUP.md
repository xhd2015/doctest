# Scenario

**Feature**: git remote -v

```
# git remote
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__URL__: 'type=string, example=git@github.com:x/y.git'\n",
		"origin\tgit@github\\.com:x/y\\.git \\(fetch\\)",
	)
	req.Actual = "origin\tgit@github.com:x/y.git (fetch)"
	return nil
}
```
