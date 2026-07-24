# Scenario

**Feature**: Same-value binding for repeated placeholders

```
# repeated __NAME__ must capture identical strings
Author -> v3 Matcher: binding state map
Matcher -> pass if same, error naming placeholder if different
```

## Steps
1. Narrow to binding match scenarios.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
