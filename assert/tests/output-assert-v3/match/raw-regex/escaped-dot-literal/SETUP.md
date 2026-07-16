# Scenario

**Feature**: V3-M2 — escaped `\.` matches literal dot only

```
# content line "a\.c" matches actual "a.c" (literal)
Author -> v3 Matcher: escaped-dot template
Matcher <- actual a.c
Matcher -> pass
```

## Steps
1. Set version-3 template body `a\.c` and actual `a.c`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", `a\.c`)
	req.Actual = "a.c"
	return nil
}
```
