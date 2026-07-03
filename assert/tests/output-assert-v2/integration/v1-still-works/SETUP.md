# Scenario

**Feature**: V2-I2 — v1 contains template still uses legacy_v1

```
# no version: 2 header — v1 <contains> block parsed by legacy_v1
Facade -> legacy_v1 Parser
Parser -> ContainsBlock AST
```

## Steps
1. Override operation to parse-only; use v1 contains template.

```go
func Setup(t *testing.T, req *Request) error {
	req.Operation = "parse"
	req.Template = "<contains>\nUsage: mytool\n  build\n</contains>"
	return nil
}
```