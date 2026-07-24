# Scenario

**Feature**: V3-E1 — dual type= and regex= on same placeholder

```
# both type=string and regex=.* on __ID__ is invalid
Author -> v3 Parser: dual type+regex definition
Parser -> parse error
```

## Steps
1. Set `__ID__: type=string, regex=.*` in header.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__ID__: type=string, regex=.*\n", "id=__ID__")
	return nil
}
```
