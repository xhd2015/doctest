# Scenario

**Feature**: V3-M4 — string placeholder non-greedy mid-line with trailing content

```
# type=string → [^\n]*? leaves room for trailing literal
Author -> v3 Matcher: Hello __NAME__!
Matcher <- Hello world!
Matcher -> pass (non-greedy string + trailing !)
```

## Steps
1. Set string placeholder mid-line with trailing `!` and matching actual.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__NAME__: type=string\n", "Hello __NAME__!")
	req.Actual = "Hello world!"
	return nil
}
```
