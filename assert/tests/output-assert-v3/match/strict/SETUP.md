# Scenario

**Feature**: Strict sequential full-match policy

```
# extra actual lines beyond template are rejected
Matcher -> match error on surplus lines
```

## Steps
1. Narrow to strict policy scenarios.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
