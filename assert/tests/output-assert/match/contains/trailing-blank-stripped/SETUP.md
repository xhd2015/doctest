# Scenario

**Feature**: C10 — trailing blank line after </contains> is stripped

```
# author writes a contains template followed by a trailing empty line
Author -> Parser: template "<contains>\nfoo\n</contains>\n\n"
# P4 strips the trailing empty line, pattern becomes contains-only
Parser -> Matcher: ContainsBlock only (no trailing literal "")
# P3 skips trailing-newline policy; P1 substring matches
Matcher <- actual "bar\nfoo\n"
Matcher -> pass
```

## Steps
1. Set template to `"<contains>\nfoo\n</contains>\n\n"` (a trailing empty line after `</contains>`).
2. Set actual to `"bar\nfoo\n"` (the fragment present, ending with `\n`).

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<contains>\nfoo\n</contains>\n\n"
	req.Actual = "bar\nfoo\n"
	return nil
}
```
