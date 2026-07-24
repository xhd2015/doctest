# Scenario

**Feature**: R1 — dot progress regex + summary

```
# realistic doctest CLI output templates
Author -> Matcher: multi-construct template
Matcher <- simulated build/help output
```

## Steps
1. Set template/actual fields for R1 — dot progress regex + summary.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<regex>\n^\\.+$\n</regex>\n  (2 Run, 2 Pass, 0 Cached, 0 Fail)"
	req.Actual = "...\n  (2 Run, 2 Pass, 0 Cached, 0 Fail)"
	return nil
}
```
