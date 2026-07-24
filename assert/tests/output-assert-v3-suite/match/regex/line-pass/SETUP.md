# Scenario

**Feature**: V3S-M4 — regex line pass

```
# .*middle.*suffix regex matches variable prefix/suffix
Matcher <- XXXSome middle contentYYYsuffix content
```

## Steps
1. Set regex-intent body line and matching actual.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("", ".*Some middle content.*suffix content")
	req.Actual = "XXXSome middle contentYYYsuffix content"
	return nil
}
```