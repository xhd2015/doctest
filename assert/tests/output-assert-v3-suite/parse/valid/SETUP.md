# Scenario

**Feature**: Valid v3 parse constructs produce typed AST nodes

```
# v3 YAML header + body produce Pattern items
Author -> v3 Parser: valid template constructs
Parser -> Pattern with Placeholder, RegexLine, OmitLine, etc.
```

## Steps
1. Expect parse success (no `ExpectParseError`).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ExpectParseError = false
	return nil
}
```
