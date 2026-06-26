# Scenario

**Feature**: H1 — hint inner text matches literally

```
# parsed pattern compared to actual output
Parser -> Matcher: Pattern
Matcher <- actual CLI output
Matcher -> pass or line-numbered diff
```

## Steps
1. Set template/actual fields for H1 — hint inner text matches literally.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<hint:id>abc</hint:id>"
	req.Actual = "abc"
	return nil
}
```
