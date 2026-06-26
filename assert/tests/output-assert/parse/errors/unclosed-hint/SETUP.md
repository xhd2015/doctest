# Scenario

**Feature**: P6 — unclosed hint

```
# invalid tag syntax rejected at parse time
Author -> Parser: malformed template
Parser -> parse error (position + message)
```

## Steps
1. Set template/actual fields for P6 — unclosed hint.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "unclosed <hint:id>abc"
	return nil
}
```
