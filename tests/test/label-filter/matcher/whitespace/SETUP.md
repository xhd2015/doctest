# Scenario

**Feature**: surrounding whitespace is ignored when parsing

```
" slow && ui " matches {slow,ui}
```

## Steps

1. Evaluate trimmed expression string.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = req
	return nil
}
```