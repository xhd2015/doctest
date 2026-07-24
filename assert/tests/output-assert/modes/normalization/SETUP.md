# Scenario

**Feature**: Newline normalization policy

```
# match options alter comparison policy
Matcher <- actual (+ Contains option or CRLF normalization)
```

## Steps
1. Compare CRLF actual against LF template.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	return nil
}
```
