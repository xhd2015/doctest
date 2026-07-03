# Scenario

**Feature**: V2-M10 — literal with dot stays pattern line

```
# version 1.0 has lone dot — not regex-intent
Matcher <- exact literal version 1.0
```

## Steps
1. Set pattern line `version 1.0` and identical actual.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "version 1.0")
	req.Actual = "version 1.0"
	return nil
}
```