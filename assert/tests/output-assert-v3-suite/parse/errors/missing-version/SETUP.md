# Scenario

**Feature**: V3S-E4 — non-dialect YAML frontmatter (no placeholder keys) falls back to legacy_v1

```
# YAML frontmatter without version and without __PLACEHOLDER__ keys is not v3 dialect
Author -> Facade: ---\nfoo: 1\n--- (no version, no placeholders)
Facade -> legacy_v1 literal parse
```

## Steps
1. Set template with plain YAML frontmatter (`foo: 1`) and no body — treated as v1 literal lines.
2. Clear `ExpectParseError` — v1 fallback should parse OK.

```go
func Setup(t *testing.T, req *Request) error {
	req.ExpectParseError = false
	req.ExpectV1Fallback = true
	req.Template = "---\nfoo: 1\n---"
	return nil
}
```
