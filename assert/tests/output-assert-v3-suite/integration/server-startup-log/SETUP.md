# Scenario

**Feature**: V3S-I1 — server startup log integration

```
# PORT placeholder + omitted stack trace + colored status line
Matcher <- realistic multi-section server log
```

## Steps
1. Set template with PORT, 4-line omit, and bold green status.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template(
		"__PORT__: type=number, example=8901, a port\n",
		"Server listen on: __PORT__\n...4 lines omitted...\n<ansi-color bold green>ready</ansi-color>",
	)
	req.Actual = "Server listen on: 8901\n" +
		"goroutine 1 [running]:\n" +
		"main.init()\n" +
		"\t/server/main.go:42\n" +
		"created by main\n" +
		greenBoldWrap("ready")
	return nil
}
```