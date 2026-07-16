# Scenario

**Feature**: V3-R1 — version: 2 still routes to legacy_v2 and matches

```
# explicit version: 2 uses legacy_v2 engine (GREEN even before v3)
Author -> Facade: version-2 template
Facade -> legacy_v2 Parser/Matcher
Matcher <- Hello alice
Matcher -> pass
```

## Steps
1. Set version: 2 string-placeholder template and matching actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "match"
	req.Template = v2Template("__USER__: type=string\n", "Hello __USER__")
	req.Actual = "Hello alice"
	return nil
}
```
