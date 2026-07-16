# Scenario

**Feature**: V3S-M16 — plain text fails color assertion

```
# <ansi-color gray> requires ANSI envelope
Matcher <- plain 1 Cached without SGR codes
```

## Steps
1. Same color template as M15 with plain actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("", "<ansi-color gray>1 Cached</ansi-color>")
	req.Actual = "1 Cached"
	return nil
}
```