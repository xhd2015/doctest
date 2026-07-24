# Scenario

**Feature**: V3S-M6 — omit three middle lines

```
# ...3 lines omitted... between greeting and closing
Matcher consumes 3 arbitrary middle lines
```

## Steps
1. Set template with USER placeholder, omit marker, and closing line.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__USER__: type=string\n", "Hello __USER__\n...3 lines omitted...\nNice to meet you")
	req.Actual = "Hello bob\nstack line 1\nstack line 2\nstack line 3\nNice to meet you"
	return nil
}
```