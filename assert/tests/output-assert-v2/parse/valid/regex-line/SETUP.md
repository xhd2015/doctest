# Scenario

**Feature**: V2-P4 — regex-intent line classified as RegexLine

```
# .* signals regex intent after protected-region scan
Author -> v2 Parser: regex body line
Parser -> RegexLine AST node
```

## Steps
1. Set body line `.*Some middle content.*suffix content` (no placeholders).

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", ".*Some middle content.*suffix content")
	return nil
}
```