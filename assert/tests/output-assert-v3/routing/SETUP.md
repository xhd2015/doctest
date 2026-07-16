# Scenario

**Feature**: Version routing — legacy engines remain reachable

```
# version: 2 → legacy_v2; pure v1 tags → legacy_v1
Author -> Facade: non-v3 templates
Facade -> legacy engines unchanged
```

## Steps
1. Routing leaves set Operation and templates explicitly.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
