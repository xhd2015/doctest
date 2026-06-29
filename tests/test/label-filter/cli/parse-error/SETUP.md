# Scenario

**Feature**: invalid `--label` expression fails before running tests

```
doctest test --label 'slow &&' -> parse error, non-zero exit
```

## Steps

1. Run test against fixture mod with invalid expression.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	_ = t
	return nil
}
```