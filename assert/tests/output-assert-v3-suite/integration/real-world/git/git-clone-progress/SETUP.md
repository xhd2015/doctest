# Scenario

**Feature**: git clone

```
# git clone
```

## Steps
1. Build v3 template and simulated actual output.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__URL__: 'type=string, example=https://x/y.git'\n",
		"Cloning into 'y'\\.\\.\\.\n...2 lines omitted...\ndone\\.",
	)
	req.Actual = "Cloning into 'y'...\nremote: Enumerating objects\nReceiving objects: 100%\ndone."
	return nil
}
```
