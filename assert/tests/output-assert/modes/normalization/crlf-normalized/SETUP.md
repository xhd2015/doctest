# Scenario

**Feature**: N2 — CRLF normalized to LF

```
# match options alter comparison policy
Matcher <- actual (+ Contains option or CRLF normalization)
```

## Steps
1. Set template/actual fields for N2 — CRLF normalized to LF.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "line1\nline2"
	req.Actual = "line1\r\nline2"
	return nil
}
```
