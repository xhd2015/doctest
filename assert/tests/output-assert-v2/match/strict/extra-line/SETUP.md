# Scenario

**Feature**: V2-M8 — strict extra actual line rejected

```
# single-line template cannot match two-line actual
Matcher rejects trailing extra line
```

## Steps
1. Set single literal line template and two-line actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "only one line")
	req.Actual = "only one line\nextra line"
	return nil
}
```