# Scenario

**Feature**: V3-M9 — ansi-color inner text is literal (dots) plus placeholder

```
# <ansi-color gray>v1.0 ready</ansi-color> — dots QuoteMeta'd inside color
# __PORT__ expands outside color span
Author -> v3 Matcher: color + placeholder
Matcher <- gray-wrapped "v1.0 ready" + port number
Matcher -> pass
```

## Steps
1. Set color span with literal dots in inner text and PORT placeholder; actual uses gray ANSI wrap.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template(
		"__PORT__: type=number\n",
		"status: <ansi-color gray>v1.0 ready</ansi-color> on __PORT__",
	)
	req.Actual = "status: " + grayWrap("v1.0 ready") + " on 8080"
	return nil
}
```
