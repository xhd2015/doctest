# Scenario

**Feature**: V2-E2 — unknown placeholder type

```
# type=boolean is not a supported placeholder type
Author -> v2 Parser: unknown type in header
Parser -> parse error
```

## Steps
1. Set `__X__: type=boolean` in header.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("__X__: type=boolean\n", "value: __X__")
	return nil
}
```