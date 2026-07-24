# Scenario

**Feature**: Strict sequential match policy

```
# no extra lines, no missing lines, trailing newline agreement
Matcher rejects length mismatches and trailing newline drift
```

## Steps
1. Templates test strict full-match invariants.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```