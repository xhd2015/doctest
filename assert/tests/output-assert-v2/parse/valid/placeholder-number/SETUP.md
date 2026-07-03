# Scenario

**Feature**: V2-P1 — number placeholder in header and body

```
# __PORT__: type=number declared in YAML, used in body
Author -> v2 Parser: placeholder-number template
Parser -> Pattern with Placeholder{PORT,number}
```

## Steps
1. Set template with `__PORT__: type=number` and body line `Server listen on: __PORT__`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("__PORT__: type=number\n", "Server listen on: __PORT__")
	return nil
}
```