# Scenario

**Feature**: Integration — realistic v2 templates and v1 compatibility

```
# multi-construct v2 templates and legacy_v1 fallback
Author -> Facade: complex templates
Matcher <- simulated CLI output
```

## Steps
1. Set `req.Operation = "match"` unless leaf overrides for parse-only v1 check.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	return nil
}
```