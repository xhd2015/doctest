# Scenario

**Feature**: Invalid v3 templates fail at parse time (or non-dialect falls back to v1)

```
# invalid v3 syntax rejected; non-placeholder YAML without version uses legacy_v1
Author -> Facade: malformed or non-v3 template
Facade -> parse error or legacy_v1 literal parse
```

## Steps
1. Expect parse failure unless `ExpectV1Fallback` is set on the leaf.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ExpectParseError = true
	return nil
}
```
