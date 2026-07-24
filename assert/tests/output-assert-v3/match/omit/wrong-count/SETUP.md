# Scenario

**Feature**: V3-M8 — omit wrong line count

```
# omit expects 3 middle lines but actual has only 2
Matcher -> line count mismatch error
```

## Steps
1. Same omit template as M7 with only 2 middle lines.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__USER__: type=string\n", "Hello __USER__\n...3 lines omitted...\nNice to meet you")
	req.Actual = "Hello bob\nonly one\nonly two\nNice to meet you"
	return nil
}
```
