# Scenario

**Feature**: V3-P7 — compact k=v placeholder definition

```
# compact form: type=number, example=8901, human explanation
Author -> v3 Parser: compact YAML placeholder def
Parser -> metadata captured (example=8901)
```

## Steps
1. Set compact placeholder __PORT__: type=number, example=8901, a port.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__PORT__: type=number, example=8901, a port\n", "Server listen on: __PORT__")
	return nil
}
```
