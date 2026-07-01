# Scenario

**Feature**: C9 — leading blank line from raw string is stripped, contains-only applies

```
# author writes a natural Go raw-string template with a leading \n
Author -> Parser: template "\n<contains>\nnot a git repository\n</contains>"
# P4 strips the leading empty line, pattern becomes contains-only
Parser -> Matcher: ContainsBlock only (no leading literal "")
# P3 skips trailing-newline policy; P1 substring matches mid-line
Matcher <- actual "wrk: /p is not a git repository\n"
Matcher -> pass
```

## Steps
1. Set template to the natural leading-newline form `"\n<contains>\nnot a git repository\n</contains>"`.
2. Set actual to a line that contains the fragment mid-line, ending with `\n`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "\n<contains>\nnot a git repository\n</contains>"
	req.Actual = "wrk: /p is not a git repository\n"
	return nil
}
```
