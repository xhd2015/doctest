# Scenario

**Feature**: P5 — bare id tag rejected

```
# invalid tag syntax rejected at parse time
Author -> Parser: malformed template
Parser -> parse error (position + message)
```

## Steps
1. Set template/actual fields for P5 — bare id tag rejected.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "<id>abc</id>"
	return nil
}
```
