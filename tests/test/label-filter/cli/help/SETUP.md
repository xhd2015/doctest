# Scenario

**Feature**: help text documents the `--label` flag

```
doctest test --help -> stdout mentions --label
```

## Steps

1. Run `doctest test --help`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	_ = t
	return nil
}
```