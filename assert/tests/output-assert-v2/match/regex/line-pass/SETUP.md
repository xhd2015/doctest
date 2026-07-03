# Scenario

**Feature**: V2-M4 — regex line pass

```
# .*middle.*suffix regex matches variable prefix/suffix
Matcher <- XXXSome middle contentYYYsuffix content
```

## Steps
1. Set regex-intent body line and matching actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", ".*Some middle content.*suffix content")
	req.Actual = "XXXSome middle contentYYYsuffix content"
	return nil
}
```