# Scenario

**Feature**: Match modes and newline normalization

```
# match options alter comparison policy
Matcher <- actual (+ Contains option or CRLF normalization)
```

## Steps
1. Leaves set match options.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Operation = "match"
	return nil
}
```
