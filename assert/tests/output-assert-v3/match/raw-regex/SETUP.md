# Scenario

**Feature**: Content lines are raw Go regex (no hasRegexIntent)

```
# unescaped metachars have RE meaning; author escapes for literals
Author -> v3 Matcher: raw regex content lines
Matcher <- actual text
```

## Steps
1. Narrow to raw-regex match scenarios.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
