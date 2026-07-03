# Scenario

**Feature**: Valid v2 parse constructs produce typed AST nodes

```
# v2 YAML header + body produce Pattern items
Author -> v2 Parser: valid template constructs
Parser -> Pattern with Placeholder, RegexLine, OmitLine, etc.
```

## Steps
1. Expect parse success (no `ExpectParseError`).

```go
func Setup(t *testing.T, req *Request) error {
	req.ExpectParseError = false
	return nil
}
```