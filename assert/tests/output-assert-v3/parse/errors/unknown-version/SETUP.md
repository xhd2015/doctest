# Scenario

**Feature**: V3-E6 — unknown version value rejected

```
# version: 9 is not a supported dialect
Author -> Facade.Parse: unknown version header
Facade -> parse error (not silent v1 fallback)
```

## Steps
1. Set YAML header `version: 9` with a simple body line.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = "---\nversion: 9\n---\nhello"
	return nil
}
```
