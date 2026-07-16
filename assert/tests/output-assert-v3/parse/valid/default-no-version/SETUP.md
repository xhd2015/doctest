# Scenario

**Feature**: V3-P2 — YAML dialect without version key defaults to v3

```
# placeholders in YAML header, no version key → v3
Author -> Facade.Parse: dialect-no-version template
Facade -> v3 Parser (default)
Parser -> Pattern with USER string placeholder
```

## Steps
1. Set template with __USER__ but no version key in YAML header.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3TemplateNoVersion("__USER__: type=string\n", "Hello __USER__")
	return nil
}
```
