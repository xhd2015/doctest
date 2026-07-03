# Scenario

**Feature**: V2-E4 — missing version:2 falls back to legacy_v1

```
# YAML frontmatter without version: 2 is not v2
Author -> Facade: ---\nfoo: 1\n--- (no version)
Facade -> legacy_v1 literal parse
```

## Steps
1. Set template without `version: 2` — treated as v1 literal lines.
2. Clear `ExpectParseError` — v1 fallback should parse OK.

```go
func Setup(t *testing.T, req *Request) error {
	req.ExpectParseError = false
	req.ExpectV1Fallback = true
	req.Template = "---\nfoo: 1\n---"
	return nil
}
```