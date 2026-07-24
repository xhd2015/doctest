# Scenario

**Feature**: V3-M6 — repeated __ID__ with different values fails

```
# two __ID__ occurrences capture different strings → match error
Author -> v3 Matcher: binding same-fail
Matcher <- id=abc again=xyz
Matcher -> match error naming ID
```

## Steps
1. Set two `__ID__` references and actual with different values.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__ID__: type=string\n", "id=__ID__\nagain=__ID__")
	req.Actual = "id=abc\nagain=xyz"
	return nil
}
```
