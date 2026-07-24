# Scenario

**Feature**: V3-E4 — unknown placeholder type without regex=

```
# type=boolean is not a supported placeholder type
Author -> v3 Parser: unknown type in header
Parser -> parse error
```

## Steps
1. Set `__X__: type=boolean` in header (no regex=).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__X__: type=boolean\n", "value: __X__")
	return nil
}
```
