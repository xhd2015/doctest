# Scenario

**Feature**: V3S-M8 — strict extra actual line rejected

```
# single-line template cannot match two-line actual
Matcher rejects trailing extra line
```

## Steps
1. Set single literal line template and two-line actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "only one line")
	req.Actual = "only one line\nextra line"
	return nil
}
```