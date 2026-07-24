# Scenario

**Feature**: V3S-E2 — unknown placeholder type

```
# type=boolean is not a supported placeholder type
Author -> v3 Parser: unknown type in header
Parser -> parse error
```

## Steps
1. Set `__X__: type=boolean` in header.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__X__: type=boolean\n", "value: __X__")
	return nil
}
```