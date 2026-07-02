# Scenario

**Feature**: C9 — explicit contains template has no leading literal blank line

```
# author writes an explicit contains template with no leading \n
Author -> Parser: template "<contains>\nnot a git repository\n</contains>"
# strict parsing preserves all template lines, so no implicit trimming occurs
Parser -> Matcher: ContainsBlock only
# contains-only matching skips trailing-newline policy; substring matches mid-line
Matcher <- actual "wrk: /p is not a git repository\n"
Matcher -> pass
```

## Steps
1. Set template to the explicit form `"<contains>\nnot a git repository\n</contains>"`.
2. Set actual to a line that contains the fragment mid-line, ending with `\n`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<contains>\nnot a git repository\n</contains>"
	req.Actual = "wrk: /p is not a git repository\n"
	return nil
}
```
