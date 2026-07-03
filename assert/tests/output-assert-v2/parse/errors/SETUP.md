# Scenario

**Feature**: Invalid v2 templates fail at parse time (or fall back to v1)

```
# invalid v2 syntax rejected; missing version:2 uses legacy_v1
Author -> Facade: malformed or non-v2 template
Facade -> parse error or legacy_v1 literal parse
```

## Steps
1. Expect parse failure unless `ExpectV1Fallback` is set on the leaf.

```go
func Setup(t *testing.T, req *Request) error {
	req.ExpectParseError = true
	return nil
}
```