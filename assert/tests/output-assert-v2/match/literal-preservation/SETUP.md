# Scenario

**Feature**: Regex detection literal-preservation edge cases

```
# lone dots, mid-line $, parens without alternation stay pattern lines
Matcher treats CLI literals as pattern, not regex
```

## Steps
1. Body lines look regex-like but lack strong intent signals.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```