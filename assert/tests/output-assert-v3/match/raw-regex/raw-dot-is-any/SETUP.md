# Scenario

**Feature**: V3-M1 — raw unescaped `.` matches any character

```
# content line "a.c" is raw RE → matches "aXc"
Author -> v3 Matcher: raw-dot template
Matcher <- actual aXc
Matcher -> pass
```

## Steps
1. Set version-3 template body `a.c` and actual `aXc`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "a.c")
	req.Actual = "aXc"
	return nil
}
```
