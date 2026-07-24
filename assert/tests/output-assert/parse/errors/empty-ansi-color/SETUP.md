# Scenario

**Feature**: P13 / AC7 — empty ansi-color inner

```
# invalid tag syntax rejected at parse time
Author -> Parser: malformed template
Parser -> parse error (position + message)
```

## Steps
1. Set template/actual fields for P13 / AC7 — empty ansi-color inner.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<ansi-color red></ansi-color>"
	return nil
}
```
