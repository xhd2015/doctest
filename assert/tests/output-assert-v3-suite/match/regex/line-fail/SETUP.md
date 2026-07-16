# Scenario

**Feature**: V3S-M5 — regex line fail

```
# regex line must fully match one actual line
Matcher <- actual with no match
```

## Steps
1. Same regex template as M4 with non-matching actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", ".*Some middle content.*suffix content")
	req.Actual = "no match here"
	return nil
}
```