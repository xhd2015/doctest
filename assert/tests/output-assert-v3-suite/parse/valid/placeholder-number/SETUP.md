# Scenario

**Feature**: V3S-P1 — number placeholder in header and body

```
# __PORT__: type=number declared in YAML, used in body
Author -> v3 Parser: placeholder-number template
Parser -> Pattern with Placeholder{PORT,number}
```

## Steps
1. Set template with `__PORT__: type=number` and body line `Server listen on: __PORT__`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__PORT__: type=number\n", "Server listen on: __PORT__")
	return nil
}
```