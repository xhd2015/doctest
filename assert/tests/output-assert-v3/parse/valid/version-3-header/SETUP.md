# Scenario

**Feature**: V3-P1 — explicit version: 3 header routes to v3

```
# version: 3 + string placeholder
Author -> Facade.Parse: version-3 template
Facade -> v3 Parser
Parser -> Pattern with USER string placeholder
```

## Steps
1. Set template with explicit version: 3 and __USER__ string placeholder.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("__USER__: type=string\n", "Hello __USER__")
	return nil
}
```
