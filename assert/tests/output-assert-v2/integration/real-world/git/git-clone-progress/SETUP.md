# Scenario

**Feature**: git clone

```
# git clone
```

## Steps
1. Build v2 template and simulated actual output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template(
		"__URL__: 'type=string, example=https://x/y.git'\n",
		"Cloning into 'y'...\n...2 lines omitted...\ndone.",
	)
	req.Actual = "Cloning into 'y'...\nremote: Enumerating objects\nReceiving objects: 100%\ndone."
	return nil
}
```
