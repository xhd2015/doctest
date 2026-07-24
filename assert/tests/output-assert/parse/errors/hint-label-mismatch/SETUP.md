# Scenario

**Feature**: P10 — hint label mismatch

```
# invalid tag syntax rejected at parse time
Author -> Parser: malformed template
Parser -> parse error (position + message)
```

## Steps
1. Set template/actual fields for P10 — hint label mismatch.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "<hint:id>abc</hint:wrong>"
	return nil
}
```
