# Scenario

**Feature**: V3-P4 — number placeholder

```
# __PORT__: type=number in header
Author -> v3 Parser: placeholder-number template
Parser -> Pattern with Placeholder{PORT,number}
```

## Steps
1. Set template with __PORT__: type=number and body Server on __PORT__.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v3Template("__PORT__: type=number\n", "Server on __PORT__")
	return nil
}
```
