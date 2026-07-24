# Scenario

**Feature**: R4 — help text contains block

```
# realistic doctest CLI output templates
Author -> Matcher: multi-construct template
Matcher <- simulated build/help output
```

## Steps
1. Set template/actual fields for R4 — help text contains block.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<contains>\nUsage: doctest\n<start-with>\n  agent\n</start-with>\n<start-with>\n  build\n</start-with>\n</contains>"
	req.Actual = "Usage: doctest [command]\n\nCommands:\n  agent    Agent commands\n  build    Build test binaries\n  test     Run tests"
	return nil
}
```
