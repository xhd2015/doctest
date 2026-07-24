# Scenario

**Feature**: V3S-P2 — string placeholder

```
# __USER__: type=string in header, used in greeting line
Author -> v3 Parser: placeholder-string template
Parser -> Pattern with Placeholder{USER,string}
```

## Steps
1. Set template with `__USER__: type=string` and body `Hello __USER__`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__USER__: type=string\n", "Hello __USER__")
	return nil
}
```