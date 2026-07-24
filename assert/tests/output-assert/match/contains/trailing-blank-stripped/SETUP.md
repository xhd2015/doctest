# Scenario

**Feature**: C10 — explicit contains template has no trailing literal blank line

```
# author writes a contains template without an extra blank line after </contains>
Author -> Parser: template "<contains>\nfoo\n</contains>"
# strict parsing preserves all template lines, so no implicit trimming occurs
Parser -> Matcher: ContainsBlock only
# contains-only matching skips trailing-newline policy; substring matches
Matcher <- actual "bar\nfoo\n"
Matcher -> pass
```

## Steps
1. Set template to `"<contains>\nfoo\n</contains>"` (no trailing empty line after `</contains>`).
2. Set actual to `"bar\nfoo\n"` (the fragment present, ending with `\n`).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<contains>\nfoo\n</contains>"
	req.Actual = "bar\nfoo\n"
	return nil
}
```
