# Scenario

**Feature**: V3S-P6 — content-only body under v3 header (two RegexLine items)

```
# header with only version: 3 + two content lines
Author -> v3 Parser: two content lines
Parser -> RegexLine+RegexLine
```

## Steps
1. Set body with two plain lines, no placeholders (v3 summarizes each as RegexLine).

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "hello\nworld")
	return nil
}
```
