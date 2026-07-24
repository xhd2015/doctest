# Scenario

**Feature**: V3-R2 — pure v1 tags still use legacy_v1

```
# no YAML dialect header — v1 <contains> block parsed by legacy_v1
Author -> Facade: v1 contains template
Facade -> legacy_v1 Parser
Parser -> ContainsBlock AST
```

## Steps
1. Parse-only with a pure v1 `<contains>` template (no version header).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Operation = "parse"
	req.Template = "<contains>\nUsage: mytool\n  build\n</contains>"
	return nil
}
```
