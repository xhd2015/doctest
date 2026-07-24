# Scenario

**Feature**: V3-E2 — invalid regex= fragment does not compile

```
# regex=[ is an incomplete character class
Author -> v3 Parser: invalid regex fragment
Parser -> parse error
```

## Steps
1. Set `__ID__: regex=[` in header (invalid Go RE fragment).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("__ID__: regex=[\n", "id=__ID__")
	return nil
}
```
