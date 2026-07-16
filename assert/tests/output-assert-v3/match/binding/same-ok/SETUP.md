# Scenario

**Feature**: V3-M5 — repeated __ID__ with same value binds

```
# two __ID__ occurrences capture identical string
Author -> v3 Matcher: binding same-ok
Matcher <- id=abc again=abc
Matcher -> pass
```

## Steps
1. Set two `__ID__` references and actual with the same captured value.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("__ID__: type=string\n", "id=__ID__\nagain=__ID__")
	req.Actual = "id=abc\nagain=abc"
	return nil
}
```
